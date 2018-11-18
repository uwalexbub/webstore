package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
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

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go validateServiceConcurrent(&wg)
	}

	wg.Wait()
}

func validateServiceConcurrent(wg *sync.WaitGroup) {
	name := getUniqueValue("medium") + ".txt"
	uploadPath := path.Join(UPLOAD_DIR, name)
	generateFile(uploadPath, medium)
	uploadFile(uploadPath)
	savePath := path.Join(DOWNLOAD_DIR, name)
	downloadFile(name, savePath)

	assertFilesAreEqual(uploadPath, savePath)
	wg.Done()
}

func validateService() {
	name := getUniqueValue("small") + ".txt"
	uploadPath := path.Join(UPLOAD_DIR, name)
	generateFile(uploadPath, small)
	uploadFile(uploadPath)
	savePath := path.Join(DOWNLOAD_DIR, name)
	downloadFile(name, savePath)

	assertFilesAreEqual(uploadPath, savePath)
}

func assertFilesAreEqual(uploadPath string, downloadPath string) {
	uploadBytes, err := ioutil.ReadFile(uploadPath)
	exitIfError(err)

	downloadBytes, err := ioutil.ReadFile(downloadPath)
	exitIfError(err)

	if len(uploadBytes) != len(downloadBytes) {
		log.Fatalf("Files %q and %q are not equal in size", uploadPath, downloadPath)
	}

	for i := 0; i < len(uploadBytes); i++ {
		if uploadBytes[i] != downloadBytes[i] {
			log.Fatalf("Files %q and %q are not equal in content", uploadPath, downloadPath)
		}
	}

	log.Printf("Files %q and %q are equal!\n", uploadPath, downloadPath)
}

func uploadFile(path string) {
	file, err := os.Open(path)
	exitIfError(err)
	defer file.Close()

	name := filepath.Base(path)
	url := fmt.Sprintf("%s/%s/%s", ENDPOINT, "upload", name)
	log.Printf("Uploading file to %s\n", url)

	resp, err := http.Post(url, CONTENT_TYPE, file)
	exitIfError(err)
	defer resp.Body.Close()

	log.Printf("Response status: %s\n", resp.Status)
}

func downloadFile(name string, savePath string) {
	url := fmt.Sprintf("%s/%s/%s", ENDPOINT, "download", name)
	log.Printf("Downloading file from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to download file: %s\n", err.Error())
	}

	file, err := os.Create(savePath)
	if err != nil {
		log.Fatalf("Failed to create file: %s\n", err.Error())
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("Failed to save file: %s\n", err.Error())
	}
	defer resp.Body.Close()

	log.Printf("Response status: %s\n", resp.Status)
}

// Generates file with random data of specified size and returns its name.
func generateFile(path string, size FileSize) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %s\n", err.Error())
	}

	data := generateData(int(size))
	file.Write(data)
	if err != nil {
		log.Fatalf("Failed to save file: %s\n", err.Error())
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

func exitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
