package pathtool

import (
	"os"
	"path/filepath"
	"strings"
)

// GetCurrentDirectory 获取当前工作目录，失败时回退到执行文件目录
func GetCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// Determine if the path exists
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// create folders
func CreateDir(path string) error {
	exist, err := PathExists(path)
	if err != nil {
		return err
	}
	if exist {
		return nil
	} else {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// create a file
func CreateFile(path string) error {
	dir := filepath.Dir(path)
	CreateDir(dir)
	exist, err := PathExists(path)
	if err != nil {
		return err
	}
	if exist {
		return nil
	} else {
		err := os.WriteFile(path, []byte(""), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
