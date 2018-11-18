package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const ENDPOINT = "http://localhost:8080"

func invokeServiceUpload(path string) {
	file, err := os.Open(path)
	exitIfError(err)
	defer file.Close()

	name := filepath.Base(path)
	url := fmt.Sprintf("%s/%s/%s", ENDPOINT, "upload", name)
	log.Printf("Uploading file to %s\n", url)

	resp, err := http.Post(url, CONTENT_TYPE, file)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected return status: %s\n", resp.Status)
	}
	exitIfError(err)
	defer resp.Body.Close()
}

func invokeServiceDownload(name string) []byte {
	url := fmt.Sprintf("%s/%s/%s", ENDPOINT, "download", name)
	log.Printf("Downloading file from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to call download method: %s\n", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected return status: %s\n", resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response from download method: %s\n", err.Error())
	}
	defer resp.Body.Close()

	return content
}
