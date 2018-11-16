package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/uwalexbub/webstore/util"
)

const (
	ENDPOINT     = "http://localhost:8080"
	UPLOAD_DIR   = "upload"   // where files are generated and uploaded from
	DOWNLOAD_DIR = "download" // where downloaded files are stored
	CONTENT_TYPE = "text/plain"
)

type FileSize int

const (
	small  FileSize = 1 * 1024 * 1024   // 1 MB
	medium          = 10 * 1024 * 1024  // 10 MB
	big             = 100 * 1024 * 1024 // 100 MB
)

const asciiLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	util.EnsureDirExists(UPLOAD_DIR)
	util.EnsureDirExists(DOWNLOAD_DIR)
	initRandSeed()
	validateService()
}

func validateService() {
	name := getUniqueValue("small") + ".txt"
	uploadPath := path.Join(UPLOAD_DIR, name)
	generateFile(uploadPath, small)
	uploadFile(uploadPath)
	savePath := path.Join(DOWNLOAD_DIR, name)
	downloadFile(name, savePath)
}

func printUniqueValues() {
	for i := 0; i < 20; i++ {
		v := getUniqueValue("test_")
		fmt.Printf("Unique value = %s\n", v)
	}
}

func uploadFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	defer file.Close()

	name := filepath.Base(path)
	url := fmt.Sprintf("%s/%s/%s", ENDPOINT, "upload", name)
	fmt.Printf("Uploading file to %s\n", url)
	resp, err := http.Post(url, CONTENT_TYPE, file)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)
}

func downloadFile(name string, savePath string) {
	url := fmt.Sprintf("%s/%s/%s", ENDPOINT, "download", name)
	fmt.Printf("Downloading file from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to download file: %s\n", err.Error())
		return
	}

	file, err := os.Create(savePath)
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", err.Error())
		return
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Failed to save file: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)
}

// Generates file with random data of specified size and returns its name.
func generateFile(path string, size FileSize) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", err.Error())
	}
	data := generateData(int(size))
	file.Write(data)
	if err != nil {
		fmt.Printf("Failed to save file: %s\n", err.Error())
	}

	defer file.Close()
}

// Generates byte array with random data of specified size
func generateData(size int) []byte {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = asciiLetters[rand.Intn(len(asciiLetters))]
	}

	return data
}

// This function is sourced from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func getUniqueValue(prefix string) string {
	const length = 12
	b := make([]byte, length)
	for i := range b {
		b[i] = asciiLetters[rand.Intn(len(asciiLetters))]
	}
	result := string(b)
	if prefix != "" {
		result = fmt.Sprintf("%s-%s", prefix, result)
	}

	return result
}

func initRandSeed() {
	rand.Seed(time.Now().UnixNano())
}
