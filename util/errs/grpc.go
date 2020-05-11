package errs

import (
	"database/sql"

	"golang.org/x/xerrors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-multierror"

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

// IsGRPCError returns determined value as boolean which a error has *status.Status
func IsGRPCError(err error) bool {
	if _, ok := status.FromError(err); ok {
		return true
	}

	if me, ok := err.(*multierror.Error); ok {
		for _, e := range me.Errors {
			return IsGRPCError(e)
		}
	}

	return false
}

// multi2grpc returns determined value as *status.Status which a error has *status.Status
func multi2grpc(err error) (*status.Status, bool) {
	if s, ok := status.FromError(err); ok {
		return s, true
	}

	if me, ok := err.(*multierror.Error); ok {
		for _, e := range me.Errors {
			return multi2grpc(e)
		}
	}

	return nil, false
}

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
	if xerrors.Is(err, ErrHTTP503) {
		return status.Error(codes.Unavailable, err.Error())
	}
	if xerrors.Is(err, ErrHTTP502) {
		return status.Error(codes.Unavailable, err.Error())
	}
	if xerrors.Is(err, ErrHTTP500) {
		return status.Error(codes.Internal, err.Error())
	}
	// if xerrors.Is(err, ErrHTTP405) {
	// 	return status.Error(codes.InvalidArgument, err.Error())
	// }
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

	if s, ok := multi2grpc(err); ok {
		return status.Error(s.Code(), err.Error())
	}

	return status.Error(codes.Unknown, err.Error())
}
