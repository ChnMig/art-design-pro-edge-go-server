package user

import (
	"strings"

	"api-server/config"
	systemuser "api-server/db/rdb/systemUser"
)

type CacheFilter struct {
	TenantID uint
	Username string
	Name     string
}

func GetUserFromCache(tenantID, id uint) (systemuser.UserCacheInfo, error) {
	if tenantID == 0 {
		return systemuser.UserCacheInfo{}, ErrInvalidTenantContext
	}
	item, err := systemuser.GetUserFromCache(id)
	if err != nil {
		return systemuser.UserCacheInfo{}, err
	}
	if item == nil {
		return systemuser.UserCacheInfo{}, nil
	}
	if item.TenantID == 0 {
		return systemuser.UserCacheInfo{}, ErrUnsafeCacheData
	}
	if item.TenantID != tenantID {
		return systemuser.UserCacheInfo{}, ErrUserNotInTenant
	}
	return *item, nil
}

func ListUsersFromCache(filter CacheFilter, page, pageSize int) ([]systemuser.UserCacheInfo, int, error) {
	if filter.TenantID == 0 {
		return nil, 0, ErrInvalidTenantContext
	}
	userList, err := systemuser.GetAllUsersFromCache()
	if err != nil {
		return nil, 0, err
	}

	filteredList, err := filterCachedUsers(userList, filter)
	if err != nil {
		return nil, 0, err
	}
	return paginateCachedUsers(filteredList, page, pageSize), len(filteredList), nil
}

func filterCachedUsers(userList []systemuser.UserCacheInfo, filter CacheFilter) ([]systemuser.UserCacheInfo, error) {
	if filter.TenantID == 0 {
		return nil, ErrInvalidTenantContext
	}
	filteredList := make([]systemuser.UserCacheInfo, 0, len(userList))
	for _, user := range userList {
		if user.TenantID == 0 {
			return nil, ErrUnsafeCacheData
		}
		if user.TenantID != filter.TenantID {
			continue
		}
		if filter.Username != "" && !strings.Contains(user.Username, filter.Username) {
			continue
		}
		if filter.Name != "" && !strings.Contains(user.Name, filter.Name) {
			continue
		}
		filteredList = append(filteredList, user)
	}
	return filteredList, nil
}

func paginateCachedUsers(filteredList []systemuser.UserCacheInfo, page, pageSize int) []systemuser.UserCacheInfo {
	total := len(filteredList)
	if total == 0 {
		return []systemuser.UserCacheInfo{}
	}

	page, pageSize = normalizeCachePagination(page, pageSize)
	lastPage := (total-1)/pageSize + 1
	if page > lastPage {
		return []systemuser.UserCacheInfo{}
	}
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	return filteredList[start:end]
}

func normalizeCachePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = config.DefaultPage
	}
	if pageSize < 1 {
		pageSize = config.DefaultPageSize
	} else if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}
	return page, pageSize
}
