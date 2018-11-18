package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"

	"github.com/uwalexbub/webstore/util"
)

const DATA_DIR = "data"

var VALID_URL_PATH = regexp.MustCompile("^/(upload|download)/([a-zA-Z0-9\\.\\-]+)$")

func main() {
	util.EnsureDirExists(DATA_DIR)
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
	log.Printf("Processing upload request for file %q\n", name)

	//logRequestData(r)
	file, err := os.Create(getFilePath(name))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Printf("Received %d bytes of data\n", n)
}

func downloadHandler(w http.ResponseWriter, r *http.Request, name string) {
	log.Printf("Processing download request for file %q\n", name)
	path := getFilePath(name)
	if !util.FileExists(path) {
		log.Printf("File %q does not exist\n", path)
		http.NotFound(w, r)
		return
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
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
