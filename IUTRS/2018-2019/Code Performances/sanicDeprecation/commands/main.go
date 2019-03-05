package main

import (
	"fmt"
	"log"
	"os"
	"sanicDeprecation"
)

func usage() {
	log.Fatal(fmt.Sprintf("Usage: %s SEARCH_DIRECTORY", os.Args[0]))
}

func printPadded(left interface{}, right interface{}) {
	fmt.Printf("%-20s %s\n", left, right)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}

	searchDirectory := os.Args[1]

	fileMatchingIsDone := make(chan bool, 1)
	fileListingIsDone := make(chan bool, 1)

	fileListingChannel := make(chan string, 1<<24)
	filePoolingChannel := make(chan *sanicDeprecation.LoadedFile, 10)

	go sanicDeprecation.ListFilesRecursive(searchDirectory, fileListingChannel, fileListingIsDone)

	go func() {
		for path := range fileListingChannel {
			go sanicDeprecation.LoadFileIntoMemory(path, filePoolingChannel)
		}
	}()

	go func() {
		for fileDefs := range filePoolingChannel {
			go sanicDeprecation.MatchFile(fileDefs)
		}

		fileMatchingIsDone <- true
	}()

	select {
	case <-fileListingIsDone:
		close(fileListingChannel)
		close(fileListingIsDone)
	case <-fileMatchingIsDone:
		close(filePoolingChannel)
		close(fileMatchingIsDone)
	}

	for pattern, matches := range sanicDeprecation.PATTERNS {
		if len(matches) > 0 {
			printPadded(pattern, "Failed")
			for _, pathMatched := range matches {
				printPadded("+", pathMatched)
			}
		} else {
			printPadded(pattern, "Success")
		}
	}
}
