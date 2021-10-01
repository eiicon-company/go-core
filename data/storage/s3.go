// Package storage gonna be implementation
// that stream io processing for memory performance.
//
package storage

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/xerrors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gobwas/glob"

	"github.com/eiicon-company/go-core/util"
	"github.com/eiicon-company/go-core/util/dsn"
	"github.com/eiicon-company/go-core/util/logger"
)

// s3Storage provides implementation s3 resource interface.
type s3Storage struct {
	Env util.Environment
	dsn *dsn.S3DSN
}

// Write will create file into the s3.
func (adp *s3Storage) Write(ctx context.Context, filename string, data []byte) error {
	var reader io.Reader = bytes.NewReader(data)

	if gzipPtn.MatchString(filename) {
		var writer *io.PipeWriter

		reader, writer = io.Pipe()
		go func() {
			gz := gzip.NewWriter(writer)
			if _, err := io.Copy(gz, bytes.NewReader(data)); err != nil {
				logger.E("[F] s3 gzip write: %s", err)
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
		return xerrors.Errorf("[F] s3 upload file failed: %w", err)
	}

	return nil
}

// Read returns file data from the s3
func (adp *s3Storage) Read(ctx context.Context, filename string) ([]byte, error) {
	file, err := ioutil.TempFile("", "s3storage")
	if err != nil {
		return nil, xerrors.Errorf("[F] s3 read file failed: %w", err)
	}

	manager := s3manager.NewDownloader(adp.dsn.Sess)
	_, err = manager.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(adp.dsn.Bucket),
		Key:    aws.String(adp.dsn.Join(filename)),
	})
	if err != nil {
		return nil, xerrors.Errorf("[F] s3 download file failed: %w", err)
	}

	var reader io.ReadCloser = file
	defer reader.Close()

	if gzipPtn.MatchString(filename) {
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return nil, xerrors.Errorf("[F] s3 gzip read failed: %w", err)
		}
	}

	data, err := ioutil.ReadAll(reader)

	os.Remove(file.Name()) // TODO: defer
	return data, err
}

// Delete will delete file from the file systems.
func (adp *s3Storage) Delete(ctx context.Context, filename string) error {
	_, err := s3.New(adp.dsn.Sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(adp.dsn.Bucket),
		Key:    aws.String(adp.dsn.Join(filename)),
	})
	return err
}

// Merge will merge file into the s3
func (adp *s3Storage) Merge(ctx context.Context, filename string, data []byte) error {
	head, _ := adp.Read(ctx, filename)
	entire := append(head, data...)

	return adp.Write(ctx, filename, entire)
}

// Files returns filename list which is traversing with glob from s3 storage.
func (adp *s3Storage) Files(ctx context.Context, ptn string) ([]string, error) {
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
func (adp *s3Storage) URL(ctx context.Context, filename string) string {
	return adp.dsn.URL(filename)
}

// String returns a URI
func (adp *s3Storage) String(ctx context.Context, filename string) string {
	return adp.dsn.String(filename)
}

// PresignedUploadURL returns a presigned upload URI
func (adp *s3Storage) PresignedUploadURL(ctx context.Context, filename string, expire time.Duration) (string, error) {
	req, _ := s3.New(adp.dsn.Sess).PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(adp.dsn.Bucket),
		Key:    aws.String(adp.dsn.Join(filename)),
	})
	return req.Presign(expire)
}

// PresignedDownloadURL returns a presigned download URI
func (adp *s3Storage) PresignedDownloadURL(ctx context.Context, filename string, expire time.Duration) (string, error) {
	req, _ := s3.New(adp.dsn.Sess).GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(adp.dsn.Bucket),
		Key:    aws.String(adp.dsn.Join(filename)),
	})
	return req.Presign(expire)
}
