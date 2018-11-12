package main

import (
	"os"
)

func ensureDirExists(dirPath string) {
	stat, err := os.Stat(dirPath)
	if os.IsNotExist(err) || (os.IsExist(err) && !stat.IsDir()) {
		os.MkdirAll(dirPath, DEFAULT_PERMS)
	}
}

func fileExists(path string) bool {
	stat, err := os.Stat(path)
	return !os.IsNotExist(err) && !stat.IsDir()
}
