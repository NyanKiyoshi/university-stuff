package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
	"sanicDeprecation"
	"sync"
)

var opts struct {
	Limit int `short:"l" long:"limit" default:"15" name:"limit of goroutines to process the files"`
	Args struct {
		DirectoryPath   string `positional-arg-name:"directory-path"`
	} `positional-args:"true" required:"true"`
}

func usage() {
	_, _ = os.Stderr.WriteString(fmt.Sprintf(
		"Usage: %s SEARCH_DIRECTORY [-limit|-l N_GOROUTINES=15]", filepath.Base(os.Args[0]),
	))
}

func printPadded(left interface{}, right interface{}) {
	fmt.Printf("%-20s %s\n", left, right)
}

func main() {
	_, err := flags.Parse(&opts);

	if err != nil {
		usage()
		os.Exit(1)
	}

	fileListingChannel := make(chan *string, opts.Limit * 2)
	filePoolingChannel := make(chan *sanicDeprecation.LoadedFile, 10)

	maxGoRoutinesChannel := make(chan bool, opts.Limit)

	var wg sync.WaitGroup

	// Launch a task to retrieve the file list
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(fileListingChannel)

		sanicDeprecation.ListFilesRecursive(
			opts.Args.DirectoryPath, fileListingChannel)
	}()

	// Launch the task of loading files from disk to the memory
	wg.Add(1)
	go func() {
		var innerWg sync.WaitGroup
		defer wg.Done()

		for path := range fileListingChannel {
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()
				sanicDeprecation.LoadFileIntoMemory(*path, filePoolingChannel)
			}()
		}

		innerWg.Wait()
		close(filePoolingChannel)
	}()

	// Launch the task of processing the loaded files
	wg.Add(1)
	go func() {
		defer wg.Done()


		for fileDefs := range filePoolingChannel {
			wg.Add(1)
			go func() {
				defer wg.Done()
				maxGoRoutinesChannel <- true
				sanicDeprecation.MatchFile(fileDefs)
				<-maxGoRoutinesChannel
			}()
		}

	}()

	// Wait for the proper tasks to finish, then close the proper channels.
	wg.Wait()

	// Print the results
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
