package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

func ReadYmlFile(filePath string, out interface{}) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, out)
}

