package cron

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"

	systemuser "api-server/db/rdb/systemUser"
)

// InitUserCacheJob 初始化用户缓存定时任务
func InitUserCacheJob() {
	// 立即执行一次缓存
	if err := systemuser.CacheAllUsers(); err != nil {
		zap.L().Error("初始化用户缓存失败", zap.Error(err))
	} else {
		zap.L().Info("初始化用户缓存成功")
	}

	// 每10分钟执行一次用户信息缓存更新
	job, err := scheduler.NewJob(
		gocron.DurationJob(
			10*60*time.Second, // 10分钟
		),
		gocron.NewTask(
			func() {
				zap.L().Info("开始执行用户缓存定时更新")
				if err := systemuser.CacheAllUsers(); err != nil {
					zap.L().Error("更新用户缓存失败", zap.Error(err))
				} else {
					zap.L().Info("更新用户缓存成功")
				}
			},
		),
	)

	if err != nil {
		zap.L().Error("创建用户缓存定时任务失败", zap.Error(err))
	} else {
		zap.L().Info("用户缓存定时任务已创建，每10分钟执行一次", zap.String("jobID", job.ID().String()))
	}
}
