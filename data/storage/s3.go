// Package storage gonna be implementation
// that stream io processing for memory performance.
//
package storage

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/eiicon-company/go-core/util"
	"github.com/eiicon-company/go-core/util/dsn"
	"github.com/eiicon-company/go-core/util/logger"
)

// s3Storage provides implementation s3 resource interface.
type s3Storage struct {
	Env util.Environment `inject:""`
	dsn *dsn.S3DSN
}

// Write will create file into the s3.
func (adp *s3Storage) Write(filename string, data []byte) error {
	var reader io.Reader = bytes.NewReader(data)

	if gzipPtn.MatchString(filename) {
		var writer *io.PipeWriter

		reader, writer = io.Pipe()
		go func() {
			gz := gzip.NewWriter(writer)
			if _, err := io.Copy(gz, bytes.NewReader(data)); err != nil {
				logger.E("s3 write: %s", err)
			}

			gz.Close()
			writer.Close()
		}()
	}

	manager := s3manager.NewUploader(adp.dsn.Sess)
	_, err := manager.Upload(&s3manager.UploadInput{
		Bucket: aws.String(adp.dsn.Bucket),
		Key:    aws.String(adp.dsn.Join(filename)),
		ACL:    aws.String(adp.dsn.ACL),
		Body:   reader,
	})

	if err != nil {
		return errors.Wrap(err, "[F] s3 upload file failed")
	}

	return nil
}

// Read returns file data from the s3
func (adp *s3Storage) Read(filename string) ([]byte, error) {
	file, err := ioutil.TempFile("", "s3storage")
	if err != nil {
		return nil, errors.Wrap(err, "[F] s3 read file failed")
	}

	manager := s3manager.NewDownloader(adp.dsn.Sess)
	_, err = manager.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(adp.dsn.Bucket),
		Key:    aws.String(adp.dsn.Join(filename)),
	})
	if err != nil {
		return nil, errors.Wrap(err, "[F] s3 download file failed")
	}

	var reader io.ReadCloser = file

	if gzipPtn.MatchString(filename) {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, errors.Wrap(err, "[F] gzip read failed")
		}
	}

	data, err := ioutil.ReadAll(reader)

	// XXX: will be fiexed to many file open
	os.Remove(file.Name())
	reader.Close()

	return data, err
}

// Merge will merge file into the s3
func (adp *s3Storage) Merge(filename string, data []byte) error {
	head, _ := adp.Read(filename)
	entire := append(head, data...)

	return adp.Write(filename, entire)
}

// Files returns filename list which is traversing with glob from s3 storage.
func (adp *s3Storage) Files(ptn string) ([]string, error) {
	g, err := glob.Compile(strings.TrimLeft(adp.dsn.Join(ptn), "/"))
	if err != nil {
		return []string{}, err
	}

	i, files := 0, []string{}
	err = s3.New(adp.dsn.Sess).ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(adp.dsn.Bucket),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		i++

		for _, obj := range p.Contents {
			if g.Match(*obj.Key) {
				files = append(files, fmt.Sprintf("s3://%s/%s", *p.Name, *obj.Key))
			}
		}

		return true
	})
	if err != nil {
		logger.Printf("Failed to retrieve list objects %s", err)
		return []string{}, err
	}

	return files, nil
}

// URL returns Public URL
func (adp *s3Storage) URL(filename string) string {
	return adp.dsn.URL(filename)
}
