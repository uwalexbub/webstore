package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/uwalexbub/webstore/util"
)

const CONTENT_TYPE = "text/plain"

var FILE_SIZES = map[string]int{
	"small":  1 * 1024 * 1024,   // 1 MB
	"medium": 10 * 1024 * 1024,  // 10 MB
	"large":  100 * 1024 * 1024, // 100 MB
}

const asciiLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	initRandSeed()
	parseCmd()
}

func parseCmd() {
	// The first element in os.Args array is always the executable itself
	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	cmd := os.Args[1]
	remainingArgs := os.Args[2:]
	if cmd == "genfiles" {
		cmdGenerateAllFiles(remainingArgs)
	} else if cmd == "functest" {
		cmdRunFunctionalTest()
	} else if cmd == "loadtest" {
		cmdRunLoadTests(remainingArgs)
	} else {
		printUsageAndExit()
	}
}

func printUsageAndExit() {
	fmt.Println(`Please specify one of the following arguments:
genfiles: generates test files
functest: runs a single functional test to validate the webservice
loadtest: runs load tests against the webservice`)
	os.Exit(1)
}

func cmdGenerateAllFiles(args []string) {
	if len(args) != 3 {
		fmt.Println(`Unrecognized arguments. Usage of gen cmd is:
genfiles <dirpath> <size> <count>

<dirpath>: path to directory where files will be created.
<size>:    label of size of files to be created. Valid values are 'small', 'medium', 'large'.
<count>:   how many files to create. `)
		os.Exit(1)
	}
	dirPath := args[0]
	sizeLabel := args[1]
	size := FILE_SIZES[sizeLabel]
	count, _ := strconv.Atoi(args[2])

	log.Println("Generating files...")
	util.EnsureDirExists(dirPath)
	for i := 0; i < count; i++ {
		name := getUniqueValue(sizeLabel) + ".txt"
		filePath := path.Join(dirPath, name)
		log.Printf("Generating file %q with %d bytes of random data", filePath, size)
		generateFile(filePath, size)
	}
}

func cmdRunFunctionalTest() {
	log.Println("Running a functional test...")
	validateService()
}

func cmdRunLoadTests(args []string) {
	if len(args) != 1 {
		fmt.Println(`Unrecognized arguments. Usage of loadtest cmd is:
laodtest <dirname>
where <dirname> is name of directory containing files for load tests`)
		os.Exit(1)
	}
	dirPath := args[0]
	log.Println("Running load tests...")

	wg := sync.WaitGroup{}
	files, err := ioutil.ReadDir(dirPath)
	exitIfError(err)

	for _, fileInfo := range files {
		wg.Add(1)
		path := filepath.Join(dirPath, fileInfo.Name())
		go runSingleTestAsync(path, &wg)
	}

	wg.Wait()
}

func runSingleTestAsync(path string, wg *sync.WaitGroup) {
	runSingleTest(path)
	wg.Done()
}

func runSingleTest(path string) {
	invokeServiceUpload(path)
	name := filepath.Base(path)
	actualContent := invokeServiceDownload(name)

	assertFileContent(path, actualContent)
}

func validateService() {
	dirPath := "tmp"
	util.EnsureDirExists(dirPath)
	sizeLabel := "small"
	name := getUniqueValue(sizeLabel) + ".txt"
	filePath := path.Join(dirPath, name)
	generateFile(filePath, FILE_SIZES[sizeLabel])

	runSingleTest(filePath)
}

func assertFileContent(expectedFilePath string, actual []byte) {
	expected, err := ioutil.ReadFile(expectedFilePath)
	exitIfError(err)

	if !util.AssertArraysAreEqual(expected, actual) {
		log.Fatalf("Content returned from webservice does not match content of previously uploaded file %q", expectedFilePath)
	}
}

// Generates file with random data of specified size and returns its name.
func generateFile(path string, size int) {
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
