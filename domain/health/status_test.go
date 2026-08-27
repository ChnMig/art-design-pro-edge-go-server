package health

import (
	"context"
	"errors"
	"testing"
)

// TestGetStatus 基础单元测试：验证健康状态的核心字段
func TestGetStatus(t *testing.T) {
	status, err := getStatusWithChecks(
		func(context.Context) error { return nil },
		func(context.Context) error { return nil },
	)

	if err != nil {
		t.Fatalf("GetStatus() error = %v, want nil", err)
	}

	if status.Status != "ok" {
		t.Errorf("GetStatus().Status = %s, want 'ok'", status.Status)
	}

	if !status.Ready {
		t.Errorf("GetStatus().Ready = %v, want true", status.Ready)
	}

	if status.Uptime <= 0 {
		t.Errorf("GetStatus().Uptime = %v, want > 0", status.Uptime)
	}

	if status.Timestamp == 0 {
		t.Errorf("GetStatus().Timestamp = %d, want non-zero", status.Timestamp)
	}
}

func TestGetStatusReportsDependencyFailure(t *testing.T) {
	status, err := getStatusWithChecks(
		func(context.Context) error { return errors.New("数据库错误") },
		func(context.Context) error { return nil },
	)
	if err == nil {
		t.Fatal("依赖不可用时应返回错误")
	}
	if status.Ready || status.Status != "unavailable" {
		t.Fatalf("依赖不可用状态错误: %#v", status)
	}
}
