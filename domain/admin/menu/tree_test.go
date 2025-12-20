package menu

import (
	"reflect"
	"testing"

	commonmenu "api-server/common/menu"
	"api-server/db/pgdb/system"

	"gorm.io/gorm"
)

func TestExtractCheckedMenuIDs(t *testing.T) {
	tree := []commonmenu.MenuResponse{
		{
			ID:            1,
			HasPermission: true,
			Children: []commonmenu.MenuResponse{
				{ID: 2, HasPermission: true},
				{ID: 3, HasPermission: false},
			},
		},
		{ID: 4, HasPermission: true},
	}

	got := extractCheckedMenuIDs(tree)
	want := []uint{1, 2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("extractCheckedMenuIDs() = %v, want %v", got, want)
	}
}

func TestExtractCheckedAuthIDs(t *testing.T) {
	tree := []commonmenu.MenuResponse{
		{
			ID: 1,
			Meta: commonmenu.MenuMeta{
				AuthList: []commonmenu.MenuAuthResp{
					{ID: 11, HasPermission: true},
					{ID: 12, HasPermission: false},
				},
			},
			Children: []commonmenu.MenuResponse{
				{
					ID: 2,
					Meta: commonmenu.MenuMeta{
						AuthList: []commonmenu.MenuAuthResp{
							{ID: 21, HasPermission: true},
						},
					},
				},
			},
		},
	}

	got := extractCheckedAuthIDs(tree)
	want := []uint{11, 21}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("extractCheckedAuthIDs() = %v, want %v", got, want)
	}
}

func TestValidateMenuScope(t *testing.T) {
	t.Run("allowed empty and menu empty", func(t *testing.T) {
		if ok := validateMenuScope(nil, nil); !ok {
			t.Fatalf("validateMenuScope() = %v, want true", ok)
		}
	})

	t.Run("allowed empty and menu non-empty", func(t *testing.T) {
		tree := []commonmenu.MenuResponse{{ID: 1}}
		if ok := validateMenuScope(tree, nil); ok {
			t.Fatalf("validateMenuScope() = %v, want false", ok)
		}
	})

	t.Run("allowed contains all", func(t *testing.T) {
		tree := []commonmenu.MenuResponse{{ID: 1, Children: []commonmenu.MenuResponse{{ID: 2}}}}
		if ok := validateMenuScope(tree, []uint{1, 2, 3}); !ok {
			t.Fatalf("validateMenuScope() = %v, want true", ok)
		}
	})

	t.Run("allowed missing child", func(t *testing.T) {
		tree := []commonmenu.MenuResponse{{ID: 1, Children: []commonmenu.MenuResponse{{ID: 2}}}}
		if ok := validateMenuScope(tree, []uint{1}); ok {
			t.Fatalf("validateMenuScope() = %v, want false", ok)
		}
	})
}

func TestValidateAuthScope(t *testing.T) {
	t.Run("allowed empty and no checked auth", func(t *testing.T) {
		tree := []commonmenu.MenuResponse{
			{
				ID: 1,
				Meta: commonmenu.MenuMeta{
					AuthList: []commonmenu.MenuAuthResp{{ID: 11, HasPermission: false}},
				},
			},
		}
		if ok := validateAuthScope(tree, nil); !ok {
			t.Fatalf("validateAuthScope() = %v, want true", ok)
		}
	})

	t.Run("allowed empty and has checked auth", func(t *testing.T) {
		tree := []commonmenu.MenuResponse{
			{
				ID: 1,
				Meta: commonmenu.MenuMeta{
					AuthList: []commonmenu.MenuAuthResp{{ID: 11, HasPermission: true}},
				},
			},
		}
		if ok := validateAuthScope(tree, nil); ok {
			t.Fatalf("validateAuthScope() = %v, want false", ok)
		}
	})

	t.Run("allowed contains checked auth", func(t *testing.T) {
		tree := []commonmenu.MenuResponse{
			{
				ID: 1,
				Meta: commonmenu.MenuMeta{
					AuthList: []commonmenu.MenuAuthResp{{ID: 11, HasPermission: true}},
				},
			},
		}
		if ok := validateAuthScope(tree, []uint{11, 12}); !ok {
			t.Fatalf("validateAuthScope() = %v, want true", ok)
		}
	})

	t.Run("allowed missing checked auth", func(t *testing.T) {
		tree := []commonmenu.MenuResponse{
			{
				ID: 1,
				Meta: commonmenu.MenuMeta{
					AuthList: []commonmenu.MenuAuthResp{{ID: 11, HasPermission: true}},
				},
			},
		}
		if ok := validateAuthScope(tree, []uint{12}); ok {
			t.Fatalf("validateAuthScope() = %v, want false", ok)
		}
	})
}

func TestFilterAuthsByScope(t *testing.T) {
	allAuths := []system.SystemMenuAuth{
		{Model: gorm.Model{ID: 1}},
		{Model: gorm.Model{ID: 2}},
		{Model: gorm.Model{ID: 3}},
	}

	got := filterAuthsByScope(allAuths, []uint{2})
	if len(got) != 1 || got[0].ID != 2 {
		t.Fatalf("filterAuthsByScope() = %v, want only ID=2", got)
	}

	got = filterAuthsByScope(allAuths, nil)
	if len(got) != 0 {
		t.Fatalf("filterAuthsByScope() = %v, want empty", got)
	}
}

