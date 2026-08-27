package user

import (
	"errors"
	"math"
	"testing"

	"api-server/config"
	systemuser "api-server/db/rdb/systemUser"
)

func TestFilterCachedUsersIsolatesTenant(t *testing.T) {
	users := []systemuser.UserCacheInfo{
		{ID: 1, TenantID: 1, Username: "alice", Name: "甲"},
		{ID: 2, TenantID: 2, Username: "alice-two", Name: "乙"},
		{ID: 3, TenantID: 1, Username: "bob", Name: "丙"},
	}

	got, err := filterCachedUsers(users, CacheFilter{TenantID: 1, Username: "alice"})
	if err != nil {
		t.Fatalf("filterCachedUsers() error = %v", err)
	}
	if len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("租户过滤结果错误: %#v", got)
	}
}

func TestFilterCachedUsersRejectsLegacyData(t *testing.T) {
	users := []systemuser.UserCacheInfo{{ID: 1, Username: "legacy"}}
	_, err := filterCachedUsers(users, CacheFilter{TenantID: 1})
	if !errors.Is(err, ErrUnsafeCacheData) {
		t.Fatalf("filterCachedUsers() error = %v, want ErrUnsafeCacheData", err)
	}
}

func TestFilterCachedUsersRequiresTenant(t *testing.T) {
	_, err := filterCachedUsers(nil, CacheFilter{})
	if !errors.Is(err, ErrInvalidTenantContext) {
		t.Fatalf("filterCachedUsers() error = %v, want ErrInvalidTenantContext", err)
	}
}

func TestPaginateCachedUsersDefendsInvalidAndOverflowingPages(t *testing.T) {
	users := make([]systemuser.UserCacheInfo, 25)
	for index := range users {
		users[index] = systemuser.UserCacheInfo{ID: uint(index + 1), TenantID: 1}
	}

	got := paginateCachedUsers(users, -1, -1)
	if len(got) != config.DefaultPageSize || got[0].ID != 1 {
		t.Fatalf("负分页参数结果错误: len=%d, first=%d", len(got), got[0].ID)
	}

	got = paginateCachedUsers(users, math.MaxInt, config.MaxPageSize)
	if len(got) != 0 {
		t.Fatalf("超大页码应返回空列表，实际长度为 %d", len(got))
	}
}
