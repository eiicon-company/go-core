package errs

import (
	"golang.org/x/xerrors"
)

var (
	// ErrHTTPUnauthorized a user need to login
	ErrHTTPUnauthorized = xerrors.New("you need logged in")

	// ErrHTTPBadRequest invalid request
	ErrHTTPBadRequest = xerrors.New("invalid argument")

	// ErrHTTPForbidden invalid request
	ErrHTTPForbidden = xerrors.New("permission denied")
)
