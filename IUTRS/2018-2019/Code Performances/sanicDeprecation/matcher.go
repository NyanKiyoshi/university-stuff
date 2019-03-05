package sanicDeprecation

import (
	"strings"
	"sync"
)

// patternsMapping defines the type of a pattern mapping against paths.
type patternsMapping map[string][]string

// patternsMappingMutex is the patterns matcher mutex.
var rxmutex sync.RWMutex

// PATTERNS defines the pattern mapping against matched files.
var PATTERNS = patternsMapping{
	"0.0.0.0":          {},
	"127.0.0.1":        {},
	"255.255.255.255":  {},
	"AF_INET":          {},
	"gethostbyaddr":    {},
	"gethostbyname":    {},
	"gethostbyname_ex": {},
	"Inet4Address":     {},
	"inet_aton":        {},
	"inet_nton":        {},
	"sockaddr_in":      {},
}

// MatchFile matches a file's content against flagged patterns.
// For now it's only doing a dumb multy-pass processing. Meaning,
// it evaluates the file with each pattern. Where, instead, we would
// want to check each "word" against a pattern.
//
// Also, we make other matchers hang as we are locking only to read
// the patterns keys. We should only read the keys and lock when we
// are actually having to write into it (rare).
// We should also handle the case where a pattern (key) would go missing,
// despite it shouldn't happen.
func MatchFile(loadedFile *LoadedFile) {
	s := string(loadedFile.content)

	// FIXME: put correct concurrency
	rxmutex.Lock()
	defer rxmutex.Unlock()

	for pattern, pathlist := range PATTERNS {
		if strings.Contains(s, pattern) {
			PATTERNS[pattern] = append(pathlist, loadedFile.path)
		}
	}
}
