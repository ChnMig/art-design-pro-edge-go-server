package user

import (
	"errors"
	"testing"

	"api-server/config"
)

func TestSuggestTenantForLoginRejectsWildcardOnlyQuery(t *testing.T) {
	oldMinLength := config.TenantMinQueryLength
	config.TenantMinQueryLength = 3
	t.Cleanup(func() { config.TenantMinQueryLength = oldMinLength })

	for _, query := range []string{"%%%", "___", `\\\`, " %_ "} {
		if _, err := SuggestTenantForLogin(query, 10); !errors.Is(err, ErrTenantQueryTooShort) {
			t.Fatalf("查询 %q 返回错误 %v，期望 ErrTenantQueryTooShort", query, err)
		}
	}
}
