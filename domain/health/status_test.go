package health

import "testing"

// TestGetStatus 基础单元测试：验证健康状态的核心字段
func TestGetStatus(t *testing.T) {
	status, err := GetStatus()

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
