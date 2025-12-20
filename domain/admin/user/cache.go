package user

import (
	"strings"

	systemuser "api-server/db/rdb/systemUser"
)

type CacheFilter struct {
	Username string
	Name     string
}

func GetUserFromCache(id uint) (systemuser.UserCacheInfo, error) {
	item, err := systemuser.GetUserFromCache(id)
	if err != nil {
		return systemuser.UserCacheInfo{}, err
	}
	if item == nil {
		return systemuser.UserCacheInfo{}, nil
	}
	return *item, nil
}

func ListUsersFromCache(filter CacheFilter, page, pageSize int) ([]systemuser.UserCacheInfo, int, error) {
	userList, err := systemuser.GetAllUsersFromCache()
	if err != nil {
		return nil, 0, err
	}

	var filteredList []systemuser.UserCacheInfo
	if filter.Username != "" || filter.Name != "" {
		filteredList = make([]systemuser.UserCacheInfo, 0, len(userList))
		for _, user := range userList {
			if filter.Username != "" && !strings.Contains(user.Username, filter.Username) {
				continue
			}
			if filter.Name != "" && !strings.Contains(user.Name, filter.Name) {
				continue
			}
			filteredList = append(filteredList, user)
		}
	} else {
		filteredList = userList
	}

	total := len(filteredList)
	if total == 0 {
		return []systemuser.UserCacheInfo{}, 0, nil
	}

	start := (page - 1) * pageSize
	if start >= total {
		return []systemuser.UserCacheInfo{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	return filteredList[start:end], total, nil
}
