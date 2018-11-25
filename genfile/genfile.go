package main

import (
	"flag"
	"log"
	"os"
	"path"

	"github.com/dustin/go-humanize"
	"github.com/uwalexbub/webstore/util"
)

var (
	dirPath       = flag.String("dir", "./data", "Path to directory where to create files.")
	sizeMegabytes = flag.Int("sizeMegabytes", 1, "Size of files in megabytes.")
	total         = flag.Int("total", 1, "How many files to create.")
)

func main() {
	flag.Parse()
	util.InitRandSeed()

	sizeBytes := *sizeMegabytes * 1024 * 1024

	log.Println("Generating files...")
	util.EnsureDirExists(*dirPath)
	for i := 0; i < *total; i++ {
		name := util.GetUniqueString(8) + ".txt"
		filePath := path.Join(*dirPath, name)
		log.Printf("Generating file %q with %s of random data", filePath, humanize.Bytes(uint64(sizeBytes)))
		generateFile(filePath, sizeBytes)
	}
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
