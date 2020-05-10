package errs

import (
	"net/http"

	"golang.org/x/xerrors"
)

var (
	// ErrHTTP503 uses as 503 ServiceUnavailable
	ErrHTTP503 = xerrors.New(http.StatusText(http.StatusServiceUnavailable))
	// ErrHTTP502 uses as 502 BadGateway
	ErrHTTP502 = xerrors.New(http.StatusText(http.StatusBadGateway))
	// ErrHTTP500 uses as 500 InternalServerError which a msg was changed from 'Internal Server Error'
	ErrHTTP500 = xerrors.New("Something went wrong")
	// ErrHTTP405 uses as 405 MethodNotAllowed
	ErrHTTP405 = xerrors.New(http.StatusText(http.StatusMethodNotAllowed))
	// ErrHTTP403 uses as 403 Forbidden
	ErrHTTP403 = xerrors.New(http.StatusText(http.StatusForbidden))
	// ErrHTTP401 uses as 401 Unauthorized
	ErrHTTP401 = xerrors.New(http.StatusText(http.StatusUnauthorized))
	// ErrHTTP400 uses as 400 BadRequest
	ErrHTTP400 = xerrors.New(http.StatusText(http.StatusBadRequest))

	// ErrHTTPForbidden denied request
	ErrHTTPForbidden = xerrors.Errorf("Permission Denied: %w", ErrHTTP403)
	// ErrHTTPUnauthorized a user need to login
	ErrHTTPUnauthorized = xerrors.Errorf("You need logged in: %w", ErrHTTP401)
	// ErrHTTPBadRequest invalid request parameter
	ErrHTTPBadRequest = xerrors.Errorf("Invalid Arguments: %w", ErrHTTP400)
)
