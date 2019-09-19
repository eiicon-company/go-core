package dsn

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

// URL returns https URL
//
// Public URL
//
// 	https://$bucket.s3.ap-southeast-2.amazonaws.com/private/$federated-identityLogo.jpg?AWSAccessKeyId=$KEY&Signature=$KEY&x-amz-security-token=$TOKEN
// 	return fmt.Sprintf("https://%s%s", dsn.Bucket, aws.StringValue(dsn.Sess.Config.Region), dsn.Join(filename))
//
func (dsn *S3DSN) URL(filename string) string {
	svc := s3.New(dsn.Sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(dsn.Bucket),
		Key:    aws.String(dsn.Key),
	})

	uri, err := req.Presign(24 * 5 * time.Hour) // TODO: Public or Private URL
	if err != nil {
		return ""
	}

	return uri
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
