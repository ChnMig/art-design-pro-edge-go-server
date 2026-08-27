package system

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSystemUserJSONNeverContainsPassword(t *testing.T) {
	encoded, err := json.Marshal(SystemUser{Account: "alice", Password: "plain-password"})
	if err != nil {
		t.Fatalf("序列化用户失败: %v", err)
	}
	if strings.Contains(string(encoded), "plain-password") || strings.Contains(string(encoded), "password") {
		t.Fatalf("用户 JSON 不应包含密码字段: %s", encoded)
	}
}
