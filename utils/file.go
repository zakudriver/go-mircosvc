package utils

import (
	"os"
	"path/filepath"
	"regexp"
)

func IsExist(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	if _, err := os.Stat(abs); err == nil {
		return true
	}
	return false
}

// 安全打开文件。 没有目标文件就创建
func SafeOpenFile(path string, flag int, mode os.FileMode) (*os.File, error) {
	if !IsExist(path) {
		r, _ := regexp.Compile(`/(\w*)\.(\w*)$`)
		rp := r.ReplaceAllString(path, "")
		if err := os.MkdirAll(rp, os.ModePerm); err != nil {
			return nil, err
		}
	}

	return os.OpenFile(path, flag, mode)
}
