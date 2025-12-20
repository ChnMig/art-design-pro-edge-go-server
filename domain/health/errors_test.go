package health

import "testing"

// TestDomainErrorsNotNil 确保领域错误已正确初始化
func TestDomainErrorsNotNil(t *testing.T) {
	if ErrServiceNotReady == nil {
		t.Errorf("ErrServiceNotReady should not be nil")
	}
	if ErrServiceUnhealthy == nil {
		t.Errorf("ErrServiceUnhealthy should not be nil")
	}
}
