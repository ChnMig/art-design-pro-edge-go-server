package pathtool

import (
	"os"
	"path/filepath"
	"strings"
)

// Get the directory where the current program is located
func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
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
