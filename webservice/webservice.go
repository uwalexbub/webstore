package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
)

var DATA_DIR = "data"
var DEFAULT_PERMS os.FileMode = 0744
var VALID_URL_PATH = regexp.MustCompile("^/(upload|download)/([a-zA-Z0-9\\.]+)$")

func main() {
	ensureDirExists(DATA_DIR)
	http.HandleFunc("/upload/", makeHandler(uploadHandler))
	http.HandleFunc("/download/", makeHandler(downloadHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, err := getName(w, r)
		if err == nil {
			fn(w, r, name)
		}
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request, name string) {
	fmt.Printf("Uploading file %q\n", name)
	logRequestData(r)
	var content []byte
	_, err := r.Body.Read(content)
	if err != nil {
		fmt.Printf("Failed to read body: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received %d bytes of data\n", len(content))
	err = ioutil.WriteFile(getFilePath(name), content, DEFAULT_PERMS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request, name string) {
	fmt.Printf("Downloading file %q\n", name)
	path := getFilePath(name)
	if !fileExists(path) {
		fmt.Printf("File %q does not exist\n", path)
		http.NotFound(w, r)
		return
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(content)
}

func logRequestData(r *http.Request) {
	fmt.Println("Headers:")
	for k, v := range r.Header {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func getName(w http.ResponseWriter, r *http.Request) (string, error) {
	m := VALID_URL_PATH.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid File Name")
	}
	return m[2], nil
}

func getFilePath(name string) string {
	return path.Join(DATA_DIR, name)
}
