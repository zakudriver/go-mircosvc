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

func OpenFileSafety(pathname string, flag int, perm os.FileMode) (f *os.File, err error) {
	f, err = os.OpenFile(pathname, flag, perm)
	if err != nil {
		reg := regexp.MustCompile(`/\w*\.?\w*$`)
		dirPath := reg.ReplaceAllString(pathname, "")
		if err := os.MkdirAll(dirPath, 0774); err != nil {
			return nil, err
		}
		return OpenFileSafety(pathname, flag, perm)
	}

	return
}
