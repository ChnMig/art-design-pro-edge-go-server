package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"api-server/api"
	"api-server/config"
	"api-server/db/pgdb"
	"api-server/db/pgdb/system"
	"api-server/util/log"

	"go.uber.org/zap"
)

func migrate() error {
	err := system.Migrate(pgdb.GetClient())
	if err != nil {
		zap.L().Error("migrate failed", zap.Error(err))
		return err
	}
	zap.L().Info("migration completed successfully")
	return nil
}

func main() {
	for _, arg := range os.Args[1:] {
		if arg == "--migrate" {
			config.RunModel = config.RunModelDevValue
			log.SetLogger()
			migrate()
			return
		}
		if arg == "--dev" {
			config.RunModel = config.RunModelDevValue
		}
	}
	log.SetLogger()
	r := api.InitApi()
	r.Run(fmt.Sprintf(":%d", config.ListenPort))
	// End of monitoring
	func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
		for {
			switch <-sigs {
			case syscall.SIGTERM, syscall.SIGINT:
				zap.L().Error("i picked up a stop signal.")
				return
			}
		}
	}()
}
