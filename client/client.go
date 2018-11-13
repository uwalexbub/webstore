package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	ENDPOINT     = "http://localhost:8080"
	DATA_DIR     = "data" // where generated files will be stored
	CONTENT_TYPE = "text/plain"
)

type FileSize int

const (
	// The keyword 'iota' assigns below constants consecutive integers, starting with 0.
	// This is Go's way of defining enums.
	small  FileSize = 1 * 1024 * 1024   // 1 MB
	medium          = 10 * 1024 * 1024  // 10 MB
	big             = 100 * 1024 * 1024 // 100 MB
)

func main() {
	initRandSeed()
	uploadFile("clienttest.txt")
}

func printUniqueValues() {
	for i := 0; i < 20; i++ {
		v := getUniqueValue("test_")
		fmt.Printf("Unique value = %s\n", v)
	}
}

func uploadFile(name string) {
	fmt.Printf("Uploading file %s\n", name)
	file, err := os.Open(name)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	defer file.Close()

	url := fmt.Sprintf("%s/%s/%s", ENDPOINT, "upload", name)
	resp, err := http.Post(url, CONTENT_TYPE, file)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)
}

func downloadFile(name string) {
	// TODO: Implement
}

// Generates file with random data and returns its name
func generateFile(size FileSize) string {
	// TODO: Impelement
	return ""
}

// This function is sourced from https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func getUniqueValue(prefix string) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	const length = 8
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	result := string(b)
	if prefix != "" {
		result = prefix + result
	}

	return result
}

func initRandSeed() {
	rand.Seed(time.Now().UnixNano())
}
