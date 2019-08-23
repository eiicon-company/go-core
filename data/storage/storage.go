package storage

import "regexp"

var (
	gzipPtn = regexp.MustCompile(".gz$") // gzipPtn uses gzip file determination.
)

type (
	// Storage provides interface for writes some of kinda data.
	Storage interface {
		Write(filename string, data []byte) error
		Read(filename string) ([]byte, error)
		Merge(filename string, data []byte) error
		Files(ptn string) ([]string, error)
	}
)
