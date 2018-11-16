package util

import (
	"os"
)

const DEFAULT_PERMS os.FileMode = 0744

func EnsureDirExists(dirPath string) {
	stat, err := os.Stat(dirPath)
	if os.IsNotExist(err) || (os.IsExist(err) && !stat.IsDir()) {
		os.MkdirAll(dirPath, DEFAULT_PERMS)
	}
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	return !os.IsNotExist(err) && !stat.IsDir()
}
