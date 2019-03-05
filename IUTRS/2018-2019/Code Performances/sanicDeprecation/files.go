package sanicDeprecation

import (
	"fmt"
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
func ListFilesRecursive(sourcePath string, output chan *string) {
	err := filepath.Walk(sourcePath, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		for _, suffix := range AllowedSuffixes {
			fmt.Printf("%s\n", path)
			if strings.HasSuffix(path, suffix) {
				output <- &path
				break
			}
		}
		return nil
	})

	if err != nil {
		log.Print("Failed to read the directory: ", err)
	}
}

// LoadFileIntoMemory loads a given file path into a LoadedFile object, in memory.
// This will read the whole file and put it as raw bytes content.
func LoadFileIntoMemory(path string, output chan *LoadedFile) {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Printf("Failed to read '%s': %s", path, err)
	}

	output <- &LoadedFile{path: path, content: data}
}
