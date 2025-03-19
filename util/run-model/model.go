package runmodel

import (
	"os"

	"api-server/config"
)

// Detect the running mode of the program
// default is dev
func Detection() {
	model := os.Getenv(config.RunModelKey)
	switch model {
	case config.RunModelDevValue:
		config.RunModel = config.RunModelDevValue
	case config.RunModelRelease:
		config.RunModel = config.RunModelRelease
	default:
		config.RunModel = config.RunModelRelease
	}
}

func IsDev() bool {
	return config.RunModel == config.RunModelDevValue
}

func IsRelease() bool {
	return config.RunModel == config.RunModelRelease
}
