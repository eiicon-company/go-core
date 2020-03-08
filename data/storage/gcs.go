// Package storage gonna be implementation
// that stream io processing for memory performance.
//
package storage

import (
	"github.com/eiicon-company/go-core/util"
	"github.com/eiicon-company/go-core/util/dsn"
)

// gcsStorage provides implementation s3 resource interface.
type gcsStorage struct {
	Env util.Environment `inject:""`
	dsn *dsn.GCSDSN
}

// Write will create file into the s3.
func (adp *gcsStorage) Write(filename string, data []byte) error {
	return nil
}

// Read returns file data from the s3
func (adp *gcsStorage) Read(filename string) ([]byte, error) {
	return nil, nil
}

// Delete will delete file from the file systems.
func (adp *gcsStorage) Delete(filename string) error {
	return nil
}

// Merge will merge file into the s3
func (adp *gcsStorage) Merge(filename string, data []byte) error {
	return nil
}

// Files returns filename list which is traversing with glob from s3 storage.
func (adp *gcsStorage) Files(ptn string) ([]string, error) {
	return nil, nil
}

// URL returns Public URL
func (adp *gcsStorage) URL(filename string) string {
	return ""
}
