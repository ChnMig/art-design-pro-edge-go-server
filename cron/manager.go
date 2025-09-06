package cron

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

var scheduler gocron.Scheduler

// InitCronJobs 初始化所有定时任务
func InitCronJobs() {
	var err error
	// 创建一个带时区的调度器
	scheduler, err = gocron.NewScheduler(gocron.WithLocation(time.Local))
	if err != nil {
		zap.L().Error("创建定时任务调度器失败", zap.Error(err))
		return
	}

	// 初始化用户缓存定时任务
	InitUserCacheJob()

	// 启动调度器
	scheduler.Start()

	zap.L().Info("定时任务调度器已启动")
}
