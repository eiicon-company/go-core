package dsn

import "testing"

func TestDSN(t *testing.T) {
	t.Helper()

	// DSN enforces implements method
	type dsn interface {
		String(filename string) string
		URL(filename string) string
	}

	var _ dsn = &FileDSN{}
	var _ dsn = &GCSDSN{}
	var _ dsn = &S3DSN{}
	// var _ DSN = &MailDSN{}
	// var _ DSN = &RedisDSN{}
}
