package health

import (
	"context"
	"errors"
	"fmt"
	"time"

	"api-server/db/pgdb"
	"api-server/db/rdb"
)

// Status 领域层的健康状态实体
// 不关心具体展示格式，仅承载核心健康信息
type Status struct {
	Status    string
	Ready     bool
	Uptime    time.Duration
	Timestamp int64
}

var startTime = time.Now()

type dependencyCheck func(context.Context) error

// GetStatus 检查 PostgreSQL 与 Redis 的实际连通性并返回服务就绪状态。
func GetStatus() (Status, error) {
	return getStatusWithChecks(pgdb.Ping, rdb.Ping)
}

func getStatusWithChecks(checkPostgres, checkRedis dependencyCheck) (Status, error) {
	status := Status{
		Status:    "ok",
		Ready:     true,
		Uptime:    time.Since(startTime),
		Timestamp: time.Now().Unix(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var checkErrors []error
	if checkPostgres == nil {
		checkErrors = append(checkErrors, errors.New("PostgreSQL 检查器未配置"))
	} else if err := checkPostgres(ctx); err != nil {
		checkErrors = append(checkErrors, fmt.Errorf("PostgreSQL 不可用: %w", err))
	}
	if checkRedis == nil {
		checkErrors = append(checkErrors, errors.New("Redis 检查器未配置"))
	} else if err := checkRedis(ctx); err != nil {
		checkErrors = append(checkErrors, fmt.Errorf("Redis 不可用: %w", err))
	}
	if len(checkErrors) > 0 {
		status.Status = "unavailable"
		status.Ready = false
		return status, fmt.Errorf("依赖健康检查失败: %w", errors.Join(checkErrors...))
	}
	return status, nil
}
