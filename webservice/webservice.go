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

const SERVICE_ENDPOINT = "localhost:8080"
const DATA_DIR = "data"
const ENCRYPTION_KEY = "This is a secret"

var VALID_URL_PATH = regexp.MustCompile("^/(upload|download)/([a-zA-Z0-9\\.\\-]+)$")

type Metrics struct {
	activeUploadRequests prometheus.Gauge
	uploadDuration       prometheus.Summary
	encryptionDuration   prometheus.Summary
	writeDuration        prometheus.Summary

	activeDownloadRequests prometheus.Gauge
	downloadDuration       prometheus.Summary
	decryptionDuration     prometheus.Summary
	readDuration           prometheus.Summary
}

var m *Metrics = &Metrics{}

func main() {
	util.RemoveDir(DATA_DIR)
	util.EnsureDirExists(DATA_DIR)
	initMetrics(m)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/upload/", makeHttpHandler(uploadHttpHandler))
	http.HandleFunc("/download/", makeHttpHandler(downloadHttpHandler))
	http.HandleFunc("/clear/", clearHandler)

	log.Printf("Started webstore service, listenting on %s", SERVICE_ENDPOINT)
	log.Fatal(http.ListenAndServe(SERVICE_ENDPOINT, nil))
}

func initMetrics(m *Metrics) {
	// Metric for number of upload requests
	m.activeUploadRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "webstore_active_upload_requests",
			Help: "Nubmer of upload requests",
		})
	prometheus.MustRegister(m.activeUploadRequests)

	// Metric for capturing duration of upload requests
	m.uploadDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "webstore_upload_duration",
			Help: "Duration of upload requests in seconds",
			// Aggregate 0.5, 0.9 and 0.99 percentiles
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
	prometheus.MustRegister(m.uploadDuration)

	// Metric for capturing duration of encryption
	m.encryptionDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "webstore_encryption_duration",
			Help: "Duration of encryption in seconds",
			// Aggregate 0.5, 0.9 and 0.99 percentiles
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
	prometheus.MustRegister(m.encryptionDuration)

	// Metric for capturing duration of writing to file
	m.writeDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "webstore_write_duration",
			Help: "Duration of write to file in seconds",
			// Aggregate 0.5, 0.9 and 0.99 percentiles
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
	prometheus.MustRegister(m.writeDuration)

	// Metric for number of download requests
	m.activeDownloadRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "webstore_active_download_requests",
			Help: "Nubmer of download requests",
		})
	prometheus.MustRegister(m.activeDownloadRequests)

	m.downloadDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "webstore_download_duration",
			Help: "Duration of download requests in seconds",
			// Aggregate 0.5, 0.9 and 0.99 percentiles
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
	prometheus.MustRegister(m.downloadDuration)

	// Metric for capturing duration of decryption
	m.decryptionDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "webstore_decryption_duration",
			Help: "Duration of decryption in seconds",
			// Aggregate 0.5, 0.9 and 0.99 percentiles
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
	prometheus.MustRegister(m.decryptionDuration)

	// Metric for capturing duration of reading from file
	m.readDuration = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "webstore_read_duration",
			Help: "Duration of read from file in seconds",
			// Aggregate 0.5, 0.9 and 0.99 percentiles
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		})
	prometheus.MustRegister(m.readDuration)
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

	// Emit metrics
	timer := prometheus.NewTimer(m.uploadDuration)
	defer timer.ObserveDuration() // 'defer' will execute the statement when parent function returns. See https://tour.golang.org/flowcontrol/12 for details.
	m.activeUploadRequests.Inc()
	defer m.activeUploadRequests.Dec()

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read body of upload request %q: %s", name, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encTimer := prometheus.NewTimer(m.encryptionDuration)
	encryptedContent := util.Encrypt(content, []byte(ENCRYPTION_KEY))
	encTimer.ObserveDuration()

	writeTimer := prometheus.NewTimer(m.writeDuration)
	err = ioutil.WriteFile(getFilePath(name), encryptedContent, util.DEFAULT_PERMS)
	if err != nil {
		log.Printf("ERROR: Failed to write file for upload request %q: %s", name, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeTimer.ObserveDuration()

	log.Printf("Encrypted and saved %q\n", name)
}

func downloadHttpHandler(w http.ResponseWriter, r *http.Request, name string) {
	log.Printf("Processing download request for file %q\n", name)

	// Emit metrics
	timer := prometheus.NewTimer(m.downloadDuration)
	defer timer.ObserveDuration()
	m.activeDownloadRequests.Inc()
	defer m.activeDownloadRequests.Dec()

	path := getFilePath(name)
	if !util.FileExists(path) {
		log.Printf("WARN: File %q does not exist\n", path)
		http.NotFound(w, r)
		return
	}

	readTimer := prometheus.NewTimer(m.readDuration)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("ERROR: Failed to read file %q: %s", name, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	readTimer.ObserveDuration()

	decTimer := prometheus.NewTimer(m.decryptionDuration)
	decryptedContent := util.Decrypt(content, []byte(ENCRYPTION_KEY))
	decTimer.ObserveDuration()

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
