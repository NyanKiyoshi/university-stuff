package sanicDeprecation

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// LoadedFile contains the definition of a in-memory loaded file.
type LoadedFile struct {
	path    string
	content []byte
}

// AllowedSuffixes lists the files suffixes to include from the file listing.
var AllowedSuffixes = [...]string{".c", ".h", ".cpp"}

// ListFilesRecursive lists to a channel every path
// contained in a starting folder.
func ListFilesRecursive(sourcePath string, output chan string, done chan bool) {
	err := filepath.Walk(sourcePath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		for _, suffix := range AllowedSuffixes {
			if strings.HasSuffix(path, suffix) {
				output <- path
				break
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	done <- true
}

// LoadFileIntoMemory loads a given file path into a LoadedFile object, in memory.
// This will read the whole file and put it as raw bytes content.
func LoadFileIntoMemory(path string, output chan *LoadedFile) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	output <- &LoadedFile{path: path, content: data}
}
