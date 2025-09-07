package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	pathtool "api-server/util/path-tool"
)

// Here are some basic configurations
// These configurations are usually generic
var (
	// listen
	ListenPort = 8080 // api listen port
	// run model
	RunModelKey      = "model"
	RunModel         = ""
	RunModelDevValue = "dev"
	RunModelRelease  = "release"
	// path
	SelfName = filepath.Base(os.Args[0])      // own file name
	AbsPath  = pathtool.GetCurrentDirectory() // current directory
	// log
	LogDir        = filepath.Join(pathtool.GetCurrentDirectory(), "log")   // log directory
	LogPath       = filepath.Join(LogDir, fmt.Sprintf("%s.log", SelfName)) // self log path
	LogMaxSize    = 50                                                     // M
	LogMaxBackups = 3                                                      // backups
	LogMaxAge     = 30                                                     // days
	LogModelDev   = "dev"                                                  // dev model
)

// Configuration variables that will be loaded from YAML
var (
	// jWT
	JWTKey        string
	JWTExpiration time.Duration
	// redis
	RedisHost     string
	RedisPassword string
	// pgsql
	PgsqlDSN string
	// admin config
	AdminPassword string
	PWDSalt       string
	// rate limit config
	LoginRatePerMinute int
	LoginBurstSize     int
	GeneralRatePerSec  int
	GeneralBurstSize   int
)

// page config
var (
	DefaultPageSize = 20 // default page size
	DefaultPage     = 1  // default page
	CancelPageSize  = -1 // cancel page size
	CancelPage      = -1 // cancel page
)

func init() {
	pathtool.CreateDir(LogDir)
	// 配置校验逻辑已移至 main.go，确保 zap logger 初始化后再校验
}
