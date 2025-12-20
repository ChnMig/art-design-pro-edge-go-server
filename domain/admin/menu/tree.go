package menu

import (
	commonmenu "api-server/common/menu"
	"api-server/db/pgdb/system"
)

func extractCheckedMenuIDs(tree []commonmenu.MenuResponse) []uint {
	var result []uint
	var walk func(items []commonmenu.MenuResponse)
	walk = func(items []commonmenu.MenuResponse) {
		for _, m := range items {
			if m.HasPermission {
				result = append(result, m.ID)
			}
			if len(m.Children) > 0 {
				walk(m.Children)
			}
		}
	}
	walk(tree)
	return result
}

func extractCheckedAuthIDs(tree []commonmenu.MenuResponse) []uint {
	var result []uint
	var walk func(items []commonmenu.MenuResponse)
	walk = func(items []commonmenu.MenuResponse) {
		for _, m := range items {
			if len(m.Meta.AuthList) > 0 {
				for _, a := range m.Meta.AuthList {
					if a.HasPermission {
						result = append(result, a.ID)
					}
				}
			}
			if len(m.Children) > 0 {
				walk(m.Children)
			}
		}
	}
	walk(tree)
	return result
}

func validateMenuScope(menus []commonmenu.MenuResponse, allowed []uint) bool {
	if len(allowed) == 0 {
		return len(menus) == 0
	}
	allowedSet := make(map[uint]struct{}, len(allowed))
	for _, id := range allowed {
		allowedSet[id] = struct{}{}
	}
	var walk func(items []commonmenu.MenuResponse) bool
	walk = func(items []commonmenu.MenuResponse) bool {
		for _, m := range items {
			if _, ok := allowedSet[m.ID]; !ok {
				return false
			}
			if len(m.Children) > 0 && !walk(m.Children) {
				return false
			}
		}
		return true
	}
	return walk(menus)
}

func validateAuthScope(menus []commonmenu.MenuResponse, allowed []uint) bool {
	if len(allowed) == 0 {
		var anyChecked bool
		var walk func(items []commonmenu.MenuResponse)
		walk = func(items []commonmenu.MenuResponse) {
			for _, m := range items {
				for _, a := range m.Meta.AuthList {
					if a.HasPermission {
						anyChecked = true
						return
					}
				}
				if len(m.Children) > 0 {
					walk(m.Children)
				}
			}
		}
		walk(menus)
		return !anyChecked
	}

	allowedSet := make(map[uint]struct{}, len(allowed))
	for _, id := range allowed {
		allowedSet[id] = struct{}{}
	}
	var okAll = true
	var walk func(items []commonmenu.MenuResponse)
	walk = func(items []commonmenu.MenuResponse) {
		for _, m := range items {
			for _, a := range m.Meta.AuthList {
				if a.HasPermission {
					if _, ok := allowedSet[a.ID]; !ok {
						okAll = false
						return
					}
				}
			}
			if len(m.Children) > 0 {
				walk(m.Children)
			}
		}
	}
	walk(menus)
	return okAll
}

func filterAuthsByScope(allAuths []system.SystemMenuAuth, allowedAuthIDs []uint) []system.SystemMenuAuth {
	if len(allowedAuthIDs) == 0 {
		return []system.SystemMenuAuth{}
	}
	allowed := make(map[uint]struct{}, len(allowedAuthIDs))
	for _, id := range allowedAuthIDs {
		allowed[id] = struct{}{}
	}
	filtered := make([]system.SystemMenuAuth, 0, len(allAuths))
	for _, a := range allAuths {
		if _, ok := allowed[a.ID]; ok {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

