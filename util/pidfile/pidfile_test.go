package pidfile

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestWriteAndRemove(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "run", "api-server.pid")

	if err := Write(path, 123); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got := strings.TrimSpace(string(data)); got != "123" {
		t.Fatalf("pid 内容=%q, want %q", got, "123")
	}

	if err := Write(path, 456); err != nil {
		t.Fatalf("Write() overwrite error = %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if got := strings.TrimSpace(string(data)); got != "456" {
		t.Fatalf("pid 覆盖后内容=%q, want %q", got, "456")
	}

	if err := Remove(path); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("pid 文件应已删除，Stat() err=%v", err)
	}

	if err := Remove(path); err != nil {
		t.Fatalf("Remove() second time error = %v", err)
	}
}

func TestWrite_EmptyPath(t *testing.T) {
	if err := Write("", os.Getpid()); err == nil {
		t.Fatalf("Write() 期望返回错误")
	}
}

func TestRemove_EmptyPath(t *testing.T) {
	if err := Remove(""); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}
}

func TestWrite_WritesNewline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pid")
	pid := 789

	if err := Write(path, pid); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	got := string(data)
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("pid 文件应以换行结尾，got=%q", got)
	}
	if strings.TrimSpace(got) != strconv.Itoa(pid) {
		t.Fatalf("pid 内容=%q, want %q", strings.TrimSpace(got), strconv.Itoa(pid))
	}
}
