package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/uwalexbub/webstore/util"
)

const CONTENT_TYPE = "text/plain"

var FILE_SIZES = map[string]int{
	"small": 1 * 1024 * 1024,  // 1 MB
	"large": 10 * 1024 * 1024, // 10 MB
}

func main() {
	util.InitRandSeed()

	if len(os.Args) < 4 {
		printUsageAndExit()
	}

	// The first element in os.Args array is the executable itself.
	dirPath := os.Args[1]
	sizeLabel := os.Args[2]
	size := FILE_SIZES[sizeLabel]
	count, _ := strconv.Atoi(os.Args[3])

	log.Println("Generating files...")
	util.EnsureDirExists(dirPath)
	for i := 0; i < count; i++ {
		name := util.GetUniqueValue(sizeLabel, 12) + ".txt"
		filePath := path.Join(dirPath, name)
		log.Printf("Generating file %q with %d bytes of random data", filePath, size)
		generateFile(filePath, size)
	}
}

func printUsageAndExit() {
	fmt.Println(`Unrecognized arguments. Valid usage is:
	genfile <dirpath> <size> <count>

	<dirpath>	path to directory where files will be created.
	<size>		label of size of files to be created. Valid values are 'small' and 'large'.
	<count>		how many files to create. `)

	os.Exit(1)
}

// Generates file with random data of specified size and returns its name.
func generateFile(path string, size int) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create file: %s\n", err.Error())
	}

	data := util.GenerateData(int(size))
	file.Write(data)
	if err != nil {
		log.Fatalf("Failed to save file: %s\n", err.Error())
	}
	defer file.Close()
}
