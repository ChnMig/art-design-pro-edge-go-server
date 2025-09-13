package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"api-server/internal/app/api/router"
	"api-server/internal/pkg/config"
	"api-server/internal/pkg/scheduler"
	"api-server/internal/repository/postgres"
	"api-server/internal/app/model/system"
	"api-server/internal/pkg/logger"

	"go.uber.org/zap"
)

func migrate() error {
	err := system.Migrate(postgres.GetClient())
	if err != nil {
		zap.L().Error("migrate failed", zap.Error(err))
		return err
	}
	zap.L().Info("migration completed successfully")
	return nil
}

func main() {
	// Load configuration from environment variables and config.yaml
	// Environment variables take precedence over config.yaml for security
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		if arg == "--migrate" {
			config.RunModel = config.RunModelDevValue
			logger.SetLogger()
			config.CheckConfig(
				config.JWTKey,
				int64(config.JWTExpiration),
				config.RedisHost,
				config.RedisPassword,
				config.PgsqlDSN,
				config.AdminPassword,
				config.PWDSalt,
			)
			migrate()
			return
		}
		if arg == "--dev" {
			config.RunModel = config.RunModelDevValue
		}
	}
	logger.SetLogger()
	config.CheckConfig(
		config.JWTKey,
		int64(config.JWTExpiration),
		config.RedisHost,
		config.RedisPassword,
		config.PgsqlDSN,
		config.AdminPassword,
		config.PWDSalt,
	)

	// 初始化定时任务
	scheduler.InitCronJobs()

	r := router.InitApi()
	go r.Run(fmt.Sprintf(":%d", config.ListenPort))

	// 监听停止信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			zap.L().Info("接收到停止信号，程序即将退出", zap.String("signal", sig.String()))
			return
		}
	}
}
