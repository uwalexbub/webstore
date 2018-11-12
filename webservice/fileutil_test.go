package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

func TestFileExists(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "fileuil_test")
	if err != nil {
		log.Fatal(err)
	}
	if !fileExists(tmpFile.Name()) {
		t.Fatalf("File %q does NOT exist but expected to exist", tmpFile.Name())
	}
	os.Remove(tmpFile.Name())
	if fileExists(tmpFile.Name()) {
		t.Fatalf("File %q exist but expected to NOT exist", tmpFile.Name())
	}

	defer os.Remove(tmpFile.Name())
}

func TestEnsureDirExists(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "fileutil_test")
	if err != nil {
		log.Fatal(err)
	}
	testDir := path.Join(tmpDir, "thedir")
	_, err = os.Stat(testDir)
	if !os.IsNotExist(err) {
		t.Fatalf("Dir %q exist but expected to NOT exist", testDir)
	}

	ensureDirExists(testDir)
	_, err = os.Stat(testDir)
	if os.IsNotExist(err) {
		t.Fatalf("Dir %q does NOT exist but expected to exist", testDir)
	}

	defer os.Remove(tmpDir)
}
