package sanicDeprecation_test

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"sanicDeprecation"
	"testing"
)

func TestListFilesRecursive_FromAbsolutePath(t *testing.T) {
	tmpPath, err := ioutil.TempDir(os.TempDir(), "filesTests")
	assert.NoError(t, err, "Failed to create temporary test directory")

	defer func() {
		err = os.RemoveAll(tmpPath)
		assert.NoError(t, err, "Failed to cleanly delete the test directory")
	}()

	// Create dummy file
	dummyFilePath := path.Join(tmpPath, "hello")
	os.Create(dummyFilePath)

	// List the expected results from the listing function
	expectedCases := [...]string{tmpPath, dummyFilePath}

	// Invoke the listing function
	outputChannel := make(chan string, 1)
	doneChannel := make(chan bool, 1)
	isDone := make(chan bool, 1)
	go sanicDeprecation.ListFilesRecursive(tmpPath, outputChannel, doneChannel)

	// Check the received bigFile is the expected one
	go func() {
		for _, expectedLine := range expectedCases {
			assert.Equal(t, expectedLine, <-outputChannel)
		}
		isDone <- true
	}()

	select {
	case <-doneChannel:
		close(outputChannel)
		close(doneChannel)
	case <-isDone:
		close(isDone)
	}

	// Ensure the output channel is closed, thus the directory listing is done
	_, ok := <-outputChannel
	assert.False(t, ok, "The output channel wasn't closed, there is still bigFile")
}

func BenchmarkListFilesRecursive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Invoke the listing function
		outputChannel := make(chan string, 1)
		doneChannel := make(chan bool, 1)
		go sanicDeprecation.ListFilesRecursive("/usr/include", outputChannel, doneChannel)

		// Read the channel
		for _ = range outputChannel {
		}
	}
}
