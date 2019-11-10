package utils

import (
	"os"
	"path/filepath"
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

