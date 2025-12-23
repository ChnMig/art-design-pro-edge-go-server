package pidfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Write 将 pid 写入指定文件路径；文件存在时会覆盖。
func Write(path string, pid int) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("pid 文件路径为空")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("创建 pid 文件目录失败: %w", err)
	}

	content := strconv.Itoa(pid) + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("写入 pid 文件失败: %w", err)
	}
	return nil
}

// Remove 删除 pid 文件；文件不存在视为成功。
func Remove(path string) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
