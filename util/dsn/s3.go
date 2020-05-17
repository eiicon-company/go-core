package dsn

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/xerrors"
)

type (
	// S3DSN s3://data-bucket/path/
	// 			 s3://data-bucket/path/?url=https://exampl.ecom:80
	S3DSN struct {
		Sess   *session.Session
		Bucket string
		Key    string
		ACL    string

		PublicURL *url.URL
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
// TODO: Get no auth or authed or private or public URL
//
func (dsn *S3DSN) URL(filename string) string {
	if dsn.PublicURL != nil {
		u, _ := url.Parse(dsn.PublicURL.String())
		u.Path = path.Join(filepath.Dir(u.Path), filename)
		return u.String()
	}

	svc := s3.New(dsn.Sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(dsn.Bucket),
		Key:    aws.String(dsn.Key),
	})

	uri, err := req.Presign(24 * 5 * time.Hour) // TODO: Auth URL: Public or Private URL
	if err != nil {
		return ""
	}

	u, _ := url.Parse(uri) // TODO: Auth URL: Public or Private URL
	u.Path = path.Join(filepath.Dir(u.Path), filename)
	u.RawQuery = ""
	u.Fragment = ""
	return u.String()
}

// S3 ...
func S3(uri string) (*S3DSN, error) {
	if uri == "" {
		return nil, ef("invalid s3 dsn")
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, xerrors.Errorf("invalid s3 dsn: %w", err)
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
		msg := "invalid s3 environment variables: %w"
		return nil, xerrors.Errorf(msg, err)
	}

	pubURL, err := url.Parse(u.Query().Get("url"))
	if err != nil {
		return nil, xerrors.Errorf("invalid url='' queryString: %w", err)
	}

	dsn := &S3DSN{
		Sess:   sess,
		Bucket: u.Host,
		Key:    u.Path,
		ACL:    "private",
	}

	if pubURL.Scheme != "" && pubURL.Host != "" {
		dsn.PublicURL = pubURL
	}

	return dsn, nil
}

func awsSession() (*session.Session, error) {
	meta, err := session.NewSession()
	if err != nil {
		msg := "aws session failed creation: %w"
		return nil, xerrors.Errorf(msg, err)
	}

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
			&ec2rolecreds.EC2RoleProvider{
				Client: ec2metadata.New(meta),
			},
		})

	var sess *session.Session

	_, err = creds.Get()

	if err == nil {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config:            aws.Config{Credentials: creds},
			SharedConfigState: session.SharedConfigDisable,
		})
	} else {
		sess, err = session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		})
	}

	if err != nil {
		msg := "invalid aws environment variables: %w"
		return nil, xerrors.Errorf(msg, err)
	}
	if aws.StringValue(sess.Config.Region) == "" {
		return nil, ef("invalid aws region is blank")
	}
	if sess.Config.Credentials == nil {
		return nil, ef("invalid aws credentials is blank")
	}

	return sess, nil
}
