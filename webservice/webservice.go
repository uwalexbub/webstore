package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/uwalexbub/webstore/util"
)

const DATA_DIR = "data"

var VALID_URL_PATH = regexp.MustCompile("^/(upload|download)/([a-zA-Z0-9\\.\\-]+)$")

const ENCRYPTION_KEY = "This is a secret"

var uploadDurations = prometheus.NewSummary(
	prometheus.SummaryOpts{
		Name:       "webstore_upload_duration",
		Help:       "Upload request duration in seconds",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})

func main() {
	util.RemoveDir(DATA_DIR)
	util.EnsureDirExists(DATA_DIR)

	prometheus.MustRegister(uploadDurations)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/upload/", makeHandler(uploadHandler))
	http.HandleFunc("/download/", makeHandler(downloadHandler))
	http.HandleFunc("/clear/", clearHandler)
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

	timer := prometheus.NewTimer(uploadDurations)
	defer timer.ObserveDuration()

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

func downloadHandler(w http.ResponseWriter, r *http.Request, name string) {
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

	// Another upload request is currently overwriting same file
	// Return error instead of maintaining locks for each file
	if len(content) == 0 {
		err := fmt.Sprintf("WARN: Failed to process download request for file %q as it is currently in use", name)
		log.Println(err)
		http.Error(w, err, http.StatusInternalServerError)
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
