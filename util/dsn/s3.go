package dsn

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

type (
	// S3DSN s3://data_bucket/path/data.flac
	S3DSN struct {
		Sess   *session.Session
		Bucket string
		Key    string
		ACL    string
	}
)

// Join returns file joined string that discards
// key's basename and then combine filename.
func (dsn *S3DSN) Join(filename string) string {
	return filepath.Join(filepath.Dir(dsn.Key), filename)
}

func (dsn *S3DSN) String(filename string) string {
	return fmt.Sprintf("s3://%s%s", dsn.Bucket, dsn.Join(filename))
}

// S3 ...
func S3(uri string) (*S3DSN, error) {
	if uri == "" {
		return nil, ef("invalid s3 dsn")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrap(err, "invalid s3 dsn")
	}
	if u.Scheme != "s3" {
		return nil, ef("invalid s3 scheme: %s", u.Scheme)
	}
	if u.Host == "" {
		return nil, ef("invalid s3 bucket is blank")
	}
	if u.Path == "" {
		return nil, ef("invalid s3 key is blank")
	}

	sess, err := awsSession()
	if err != nil {
		msg := "invalid s3 environment variables"
		return nil, errors.Wrap(err, msg)
	}

	dsn := &S3DSN{
		Sess:   sess,
		Bucket: u.Host,
		Key:    u.Path,
		ACL:    "private",
	}
	return dsn, nil
}

func awsSession() (*session.Session, error) {
	meta, err := session.NewSession()
	if err != nil {
		msg := "aws session failed creation"
		return nil, errors.Wrap(err, msg)
	}

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(meta),
			},
		})

	if _, err := creds.Get(); err != nil {
		msg := "invalid aws environment variables"
		return nil, errors.Wrap(err, msg)
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Credentials: creds},
		SharedConfigState: session.SharedConfigDisable,
	})
	if err != nil {
		msg := "invalid aws environment variables"
		return nil, errors.Wrap(err, msg)
	}
	if aws.StringValue(sess.Config.Region) == "" {
		return nil, ef("invalid aws region is blank")
	}

	return sess, nil
}
