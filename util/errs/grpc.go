package errs

import (
	"database/sql"

	"golang.org/x/xerrors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eiicon-company/go-core/util/repo"
)

var (
	// ErrGRPCUnauthenticated a user need to login
	ErrGRPCUnauthenticated = status.Error(codes.Unauthenticated, "you need logged in")

	// ErrGRPCInvalidArgument invalid request
	ErrGRPCInvalidArgument = status.Error(codes.InvalidArgument, "invalid argument")

	// ErrGRPCInternal server error
	ErrGRPCInternal = status.Error(codes.Internal, "something went wrong")
)

// GRPCError returns grpc status error
func GRPCError(err error) error {
	if err == nil {
		return nil
	}

	// GRPC Status Code
	if _, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return err
	}

	// Others

	if xerrors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, err.Error())
	}
	if xerrors.Is(err, repo.ErrExists) {
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if xerrors.Is(err, ErrHTTPUnauthorized) {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	if xerrors.Is(err, ErrHTTPBadRequest) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if xerrors.Is(err, ErrHTTPForbidden) {
		return status.Error(codes.PermissionDenied, err.Error())
	}

	return status.Error(codes.Unknown, err.Error())
}
