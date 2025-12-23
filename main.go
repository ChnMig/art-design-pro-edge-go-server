package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"go.uber.org/zap"

	"api-server/api"
	"api-server/api/middleware"
	"api-server/config"
	"api-server/cron"
	"api-server/db/pgdb"
	"api-server/db/pgdb/system"
	"api-server/util/acme"
	"api-server/util/log"
	pathtool "api-server/util/path-tool"
	"api-server/util/pidfile"
	runmodel "api-server/util/run-model"
	"api-server/util/tlsfile"
)

var CLI struct {
	Dev     bool `help:"以开发模式运行" short:"d"`
	Migrate bool `help:"执行数据库迁移后退出"`
	Version bool `help:"显示版本信息" short:"v"`
}

var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func migrate() error {
	if err := system.Migrate(pgdb.GetClient()); err != nil {
		zap.L().Error("数据库迁移失败", zap.Error(err))
		return err
	}
	zap.L().Info("数据库迁移完成")
	return nil
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("api-server"),
		kong.Description("art-design-pro-edge 后端服务"),
		kong.UsageOnError(),
	)

	if CLI.Version {
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		os.Exit(0)
	}

	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		ctx.Exit(1)
	}

	if CLI.Dev {
		config.RunModel = config.RunModelDevValue
	} else {
		runmodel.Detection()
	}

	// 仅在生产模式创建日志目录，避免测试/子包初始化时散落空 log 目录
	if runmodel.IsRelease() {
		if err := pathtool.CreateDir(config.LogDir); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
			ctx.Exit(1)
		}
	}

	log.GetLogger()
	log.StartMonitor()

	config.WatchConfig(func() {
		zap.L().Info("配置重载成功",
			zap.Int("port", config.ListenPort),
			zap.Duration("jwt_expiration", config.JWTExpiration),
			zap.Bool("global_rate_limit", config.EnableRateLimit),
		)
	})

	config.CheckConfig(
		config.JWTKey,
		int64(config.JWTExpiration),
		config.RedisHost,
		config.RedisPassword,
		config.PgsqlDSN,
		config.AdminPassword,
		config.PWDSalt,
	)

	if CLI.Migrate {
		exitCode := 0
		if err := migrate(); err != nil {
			exitCode = 1
		}
		log.StopMonitor()
		ctx.Exit(exitCode)
	}

	cron.InitCronJobs()

	r := api.InitApi()
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.ListenPort),
		Handler:        r,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		IdleTimeout:    config.IdleTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}

	acmeCtx := acme.Setup(srv)
	tlsFileCtx := tlsfile.Setup(srv)

	// 监听停止信号（尽早注册，避免启动阶段收到信号时错过清理流程）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	pidFilePath := config.PidFile
	// 写入 pid 文件（存在则覆盖，确保每次启动都会刷新）
	if pidFilePath != "" {
		pid := os.Getpid()
		if err := pidfile.Write(pidFilePath, pid); err != nil {
			zap.L().Error("写入 pid 文件失败",
				zap.String("pid_file", pidFilePath),
				zap.Error(err),
			)
			log.StopMonitor()
			ctx.Exit(1)
		}
		zap.L().Info("PID 文件已写入",
			zap.String("pid_file", pidFilePath),
			zap.Int("pid", pid),
		)
	}

	if acmeCtx.Enabled && acmeCtx.HTTPServer != nil {
		go func() {
			zap.L().Info("ACME HTTP 挑战服务启动", zap.String("addr", acmeCtx.HTTPServer.Addr))
			if err := acmeCtx.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				zap.L().Error("ACME HTTP 挑战服务异常退出", zap.Error(err))
			}
		}()
	}

	serverErrCh := make(chan error, 1)
	go func() {
		var err error
		if acmeCtx.Enabled || tlsFileCtx.Enabled {
			zap.L().Info("HTTPS 服务启动中", zap.String("addr", srv.Addr))
			err = srv.ListenAndServeTLS("", "")
		} else {
			zap.L().Info("HTTP 服务启动中", zap.String("addr", srv.Addr))
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			serverErrCh <- err
		}
	}()

	exitCode := 0
	select {
	case sig := <-quit:
		zap.L().Info("接收到停止信号，开始优雅退出", zap.String("signal", sig.String()))
	case err := <-serverErrCh:
		exitCode = 1
		zap.L().Error("HTTP 服务异常退出，开始执行清理与退出",
			zap.Error(err),
		)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)

	if err := srv.Shutdown(shutdownCtx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			zap.L().Error("HTTP 服务强制退出", zap.Error(err))
		}
	}

	if acmeCtx.Enabled && acmeCtx.HTTPServer != nil {
		if err := acmeCtx.HTTPServer.Shutdown(shutdownCtx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				zap.L().Error("ACME HTTP 挑战服务关闭失败", zap.Error(err))
			}
		}
	}

	middleware.CleanupAllLimiters()

	log.StopMonitor()

	// 删除 pid 文件（文件不存在视为成功）
	if err := pidfile.Remove(pidFilePath); err != nil {
		zap.L().Warn("删除 pid 文件失败",
			zap.String("pid_file", pidFilePath),
			zap.Error(err),
		)
	}

	cancel()
	zap.L().Info("服务退出完成")
	ctx.Exit(exitCode)
}
