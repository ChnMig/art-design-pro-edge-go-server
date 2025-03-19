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

// These configurations need to be modified as needed
var (
	// jWT
	JWTKey        = "CvXPiv34e2474LC5Xj7IP" // 务必在部署前对 key 进行修改
	JWTExpiration = time.Hour * 12
)

// redis
var (
	RedisHost     = "127.0.0.1:6379"        // 修改为自己的 redis 地址
	RedisPassword = "izpXvn894uW2HFbyP5OGr" // 修改为自己的 redis 密码
)

// pgsql
var (
	PgsqlDSN = "host=127.0.0.1 user=postgres password=kL81xnDWo221FHFRX8GnP dbname=server port=5432 sslmode=disable TimeZone=Asia/Shanghai"
)

// admin config
var (
	AdminPassword = "123456"                // 管理员密码, 务必修改为自己的密码
	PWDSalt       = "rHECMvW3el1zhpdzgx9dY" // 数据库存储密码时的盐, 务必重新生成, 且不可泄露, 不可更改
)

func init() {
	pathtool.CreateDir(LogDir)
}
