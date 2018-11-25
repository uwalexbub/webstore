package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/uwalexbub/webstore/util"
)

const SERVICE_ENDPOINT = "localhost:8080"
const DATA_DIR = "data"
const ENCRYPTION_KEY = "This is a secret"

var VALID_URL_PATH = regexp.MustCompile("^/(upload|download)/([a-zA-Z0-9\\.\\-]+)$")

func main() {
	util.RemoveDir(DATA_DIR)
	util.EnsureDirExists(DATA_DIR)

	http.HandleFunc("/upload/", makeHttpHandler(uploadHttpHandler))
	http.HandleFunc("/download/", makeHttpHandler(downloadHttpHandler))
	http.HandleFunc("/clear/", clearHandler)

	log.Printf("Started webstore service, listenting on %s", SERVICE_ENDPOINT)
	log.Fatal(http.ListenAndServe(SERVICE_ENDPOINT, nil))
}

func makeHttpHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, err := getName(w, r)
		if err == nil {
			fn(w, r, name)
		}
	}
}

func uploadHttpHandler(w http.ResponseWriter, r *http.Request, name string) {
	log.Printf("Processing upload request for file %q\n", name)

	//logRequestData(r)
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read body of upload request %q: %s", name, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encryptedContent := util.Encrypt(content, []byte(ENCRYPTION_KEY))

	err = ioutil.WriteFile(getFilePath(name), encryptedContent, util.DEFAULT_PERMS)
	if err != nil {
		log.Printf("ERROR: Failed to write file for upload request %q: %s", name, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Encrypted and saved %q\n", name)
}

func downloadHttpHandler(w http.ResponseWriter, r *http.Request, name string) {
	log.Printf("Processing download request for file %q\n", name)
	path := getFilePath(name)
	if !util.FileExists(path) {
		log.Printf("WARN: File %q does not exist\n", path)
		http.NotFound(w, r)
		return
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("ERROR: Failed to read file %q: %s", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	decryptedContent := util.Decrypt(content, []byte(ENCRYPTION_KEY))
	w.Header().Set("Content-Type", "text/plain")
	w.Write(decryptedContent)
	log.Printf("Decrypted and returned %q", name)
}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	if err := util.RemoveDir(DATA_DIR); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.EnsureDirExists(DATA_DIR)
	log.Println("Cleared data directory")
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
