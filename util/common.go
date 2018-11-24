package util

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

const DEFAULT_PERMS os.FileMode = 0744
const AsciiLetters = "abcdefghijklmnopqrstuvwxyz"

func RemoveDir(dirPath string) error {
	return os.RemoveAll(dirPath)
}

func EnsureDirExists(dirPath string) {
	os.MkdirAll(dirPath, DEFAULT_PERMS)
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	return !os.IsNotExist(err) && !stat.IsDir()
}

func AssertArraysAreEqual(expected []byte, actual []byte) bool {
	if len(expected) != len(actual) {
		return false
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			return false
		}
	}

	return true
}

func InitRandSeed() {
	rand.Seed(time.Now().UnixNano())
}

func ExitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// This function is sourced from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func GetUniqueValue(prefix string, size int) string {
	const length = 12
	b := make([]byte, length)
	for i := range b {
		b[i] = AsciiLetters[rand.Intn(len(AsciiLetters))]
	}
	result := string(b)
	if prefix != "" {
		result = fmt.Sprintf("%s-%s", prefix, result)
	}

	return result
}

// Generates byte array of specified size with random data
func GenerateData(size int) []byte {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = AsciiLetters[rand.Intn(len(AsciiLetters))]
	}

	return data
}
