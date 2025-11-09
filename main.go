package main

import (
	"context"
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
	"api-server/util/log"
	runmodel "api-server/util/run-model"
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

	log.GetLogger()
	log.StartMonitor()
	defer log.StopMonitor()

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
		if err := migrate(); err != nil {
			ctx.Exit(1)
		}
		ctx.Exit(0)
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

	go func() {
		zap.L().Info("HTTP 服务启动中", zap.Int("port", config.ListenPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("HTTP 服务启动失败", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	zap.L().Info("接收到停止信号，开始优雅退出", zap.String("signal", sig.String()))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("HTTP 服务强制退出", zap.Error(err))
	}

	middleware.CleanupAllLimiters()
	zap.L().Info("服务退出完成")
	ctx.Exit(0)
}
