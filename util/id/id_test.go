package id

import (
	"regexp"
	"testing"
)

func TestIssueID(t *testing.T) {
	// 测试生成 ID
	id1 := IssueID()
	if id1 == "" {
		t.Fatal("IssueID returned empty string")
	}

	// 测试 ID 应该是数字字符串
	matched, err := regexp.MatchString(`^\d+$`, id1)
	if err != nil {
		t.Fatalf("Regex error: %v", err)
	}
	if !matched {
		t.Errorf("IssueID should return numeric string, got: %s", id1)
	}

	// 测试生成的 ID 应该是唯一的
	id2 := IssueID()
	if id1 == id2 {
		t.Error("IssueID should generate unique IDs")
	}
}

func TestIssueID_WhenFlakeNil(t *testing.T) {
	oldFlake := flake
	flake = nil
	t.Cleanup(func() {
		flake = oldFlake
	})

	id1 := IssueID()
	if id1 == "" {
		t.Fatal("IssueID returned empty string when flake is nil")
	}

	matched, err := regexp.MatchString(`^\d+$`, id1)
	if err != nil {
		t.Fatalf("Regex error: %v", err)
	}
	if !matched {
		t.Errorf("IssueID should return numeric string when flake is nil, got: %s", id1)
	}
}

func TestGenerateID(t *testing.T) {
	// 测试生成唯一 ID
	id1 := GenerateID()
	if id1 == "" {
		t.Fatal("GenerateID returned empty string")
	}

	// MD5 哈希应该是 32 个十六进制字符
	if len(id1) != 32 {
		t.Errorf("GenerateID should return 32 character MD5 hash, got length: %d", len(id1))
	}

	// 测试是否为有效的十六进制字符串
	matched, err := regexp.MatchString(`^[a-f0-9]{32}$`, id1)
	if err != nil {
		t.Fatalf("Regex error: %v", err)
	}
	if !matched {
		t.Errorf("GenerateID should return hex string, got: %s", id1)
	}

	// 测试生成的 ID 应该是唯一的
	id2 := GenerateID()
	if id1 == id2 {
		t.Error("GenerateID should generate unique IDs")
	}
}

func TestIDUniqueness(t *testing.T) {
	// 批量生成 ID 并检查唯一性
	ids := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		id := IssueID()
		if ids[id] {
			t.Errorf("Duplicate ID found: %s", id)
		}
		ids[id] = true
	}

	if len(ids) != count {
		t.Errorf("Expected %d unique IDs, got %d", count, len(ids))
	}
}

func BenchmarkIssueID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IssueID()
	}
}

func BenchmarkGenerateID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateID()
	}
}
