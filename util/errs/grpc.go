package errs

import (
	"database/sql"

	"golang.org/x/xerrors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eiicon-company/go-core/util/repo"
)

var (
	// ErrGRPCInternal server error
	ErrGRPCInternal = status.Error(codes.Internal, "Something went wrong")
	// ErrGRPCUnauthenticated a user need to login
	ErrGRPCUnauthenticated = status.Error(codes.Unauthenticated, "You need logged in")
	// ErrGRPCInvalidArgument invalid request
	ErrGRPCInvalidArgument = status.Error(codes.InvalidArgument, "Invalid Argument")
)

// GRPCError returns grpc status error
func GRPCError(err error) error {
	if err == nil {
		return nil
	}

	// Others

	if xerrors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, err.Error())
	}
	if xerrors.Is(err, repo.ErrExists) {
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if xerrors.Is(err, ErrHTTP403) {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	if xerrors.Is(err, ErrHTTP401) {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	if xerrors.Is(err, ErrHTTP400) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	// With GRPC Status Code

	if _, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return err
	}

	if xerrors.Is(err, ErrGRPCInternal) {
		return status.Error(codes.Internal, err.Error())
	}
	if xerrors.Is(err, ErrGRPCUnauthenticated) {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	if xerrors.Is(err, ErrGRPCInvalidArgument) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return status.Error(codes.Unknown, err.Error())
}
