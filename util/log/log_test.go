package log

import "testing"

func TestRedactBoundParamsRecursively(t *testing.T) {
	input := struct {
		Account       string                 `json:"account"`
		Password      string                 `json:"password"`
		AdminPassword string                 `json:"admin_password"`
		Nested        map[string]interface{} `json:"nested"`
	}{
		Account:       "alice",
		Password:      "plain-password",
		AdminPassword: "admin-password",
		Nested:        map[string]interface{}{"access_token": "token-value", "safe": "visible"},
	}

	redacted, ok := redactBoundParams(input).(map[string]interface{})
	if !ok {
		t.Fatalf("脱敏结果类型错误: %#v", redactBoundParams(input))
	}
	if redacted["password"] != redactedValue || redacted["admin_password"] != redactedValue {
		t.Fatalf("密码字段未脱敏: %#v", redacted)
	}
	if redacted["account"] != "alice" {
		t.Fatalf("普通字段不应被改写: %#v", redacted)
	}
	nested := redacted["nested"].(map[string]interface{})
	if nested["access_token"] != redactedValue || nested["safe"] != "visible" {
		t.Fatalf("嵌套字段脱敏错误: %#v", nested)
	}
}

func TestRedactStringValuesDoesNotMutateSource(t *testing.T) {
	source := map[string][]string{"authorization": []string{"Bearer secret"}, "name": []string{"alice"}}
	redacted := redactStringValues(source)
	if redacted["authorization"][0] != redactedValue {
		t.Fatalf("Authorization 未脱敏: %#v", redacted)
	}
	if source["authorization"][0] != "Bearer secret" {
		t.Fatalf("脱敏函数不应修改原始参数: %#v", source)
	}
}
