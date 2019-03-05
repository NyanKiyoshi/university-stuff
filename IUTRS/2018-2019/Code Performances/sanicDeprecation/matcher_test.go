package sanicDeprecation_test

import (
	"sanicDeprecation"
	"testing"
)

var bigFile, smallFile *sanicDeprecation.LoadedFile

func init() {
	var file = make(chan *sanicDeprecation.LoadedFile, 1)
	go sanicDeprecation.LoadFileIntoMemory("/usr/include/sqlite3.h", file)
	go sanicDeprecation.LoadFileIntoMemory("/usr/include/linux/atm.h", file)
	bigFile = <-file
	smallFile = <-file
}

func BenchmarkMatchFile_AgainstBigFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sanicDeprecation.MatchFile(bigFile)
	}
}

func BenchmarkMatchFile_AgainstSmallFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sanicDeprecation.MatchFile(smallFile)
	}
}
