package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/uwalexbub/webstore/util"
)

var (
	// help = flag.Bool("help", false, "Print usage of the program.")
	parallelism = flag.Int("parallelism", 1, "Number of parallel tests running continuously.")
	bytesMin    = flag.Int("bytes.min", 1*1024*1024, "Lower bound of test data size.")
	bytesMax    = flag.Int("bytes.max", 10*1024*1024, "Upper bound of amount of test data size.")
)

const SERVICE_ENDPOINT = "http://localhost:8080"
const CONTENT_TYPE = "text/plain"

func main() {
	flag.Parse()
	run()
}

func run() {
	util.InitRandSeed()
	invokeClear()

	dataBank := generateDataBank()
	runTests(dataBank)
}

func generateDataBank() *[]byte {
	size := 20 * (*bytesMax)
	log.Printf("Generating bank of random data of %s...\n", humanize.Bytes(uint64(size)))
	dataBank := util.GenerateData(size)
	return &dataBank
}

func runTests(dataBank *[]byte) {
	log.Printf("Starting %d test treads\n", parallelism)
	wg := sync.WaitGroup{}
	for i := 0; i < *parallelism; i++ {
		wg.Add(1)
		go runContinuousTestAsync(dataBank, &wg)
	}
	wg.Wait()
}

func runContinuousTestAsync(dataBank *[]byte, wg *sync.WaitGroup) {
	for {
		runSingleTest(dataBank, false)
		time.Sleep(time.Millisecond)
	}
	// wg.Done()
}

func runSingleTest(dataBank *[]byte, forceSuccess bool) {
	name := util.GetUniqueValue("", 12)

	start := rand.Intn(len(*dataBank) - *bytesMax)
	end := start + *bytesMin + rand.Intn(*bytesMax-*bytesMin)
	expectedBytes := (*dataBank)[start:end]

	invokeUpload(name, expectedBytes)
	actualBytes, err := invokeDownload(name)

	if err != nil {
		if forceSuccess {
			log.Fatal(err)
		}
		log.Printf("WARN: Failed to download %q: %s", name, err.Error())
	} else if err == nil {
		util.AssertArraysAreEqual(expectedBytes, actualBytes)
	}
}

func invokeUpload(name string, data []byte) {
	url := fmt.Sprintf("%s/%s/%s", SERVICE_ENDPOINT, "upload", name)
	log.Printf("Uploading %s to %s\n", humanize.Bytes(uint64(len(data))), url)

	reader := bytes.NewReader(data)
	resp, err := http.Post(url, CONTENT_TYPE, reader)
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to upload: %s\n", resp.Status)
	}
	util.ExitIfError(err)
	defer resp.Body.Close()
}

func invokeDownload(name string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", SERVICE_ENDPOINT, "download", name)
	log.Printf("Downloading from %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to call download method: %s\n", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("Unexpected return status: %s\n", resp.Status)
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response from download method: %s\n", err.Error())
	}
	defer resp.Body.Close()

	return content, nil
}

func invokeClear() {
	url := fmt.Sprintf("%s/clear", SERVICE_ENDPOINT)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to call clear method: %s\n", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Unexpected return status: %s\n", resp.Status)
	}
	resp.Body.Close()
}
