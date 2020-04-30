package dsn

import "golang.org/x/xerrors"

func ef(format string, a ...interface{}) error {
	return xerrors.Errorf(format, a...)
}

// filePublicURL Http URL
var filePublicURL = "http://localhost:8000"
