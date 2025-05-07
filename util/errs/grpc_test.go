package errs

import (
	"database/sql"
	"testing"

	"golang.org/x/xerrors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"

	"github.com/eiicon-company/go-core/util/repo"
)

func TestGRPCError(t *testing.T) {
	t.Helper()

	var err error

	err = xerrors.Errorf("ErrHTTPUnauthorized: %w", ErrHTTPUnauthorized)
	require.Equal(t, codes.Unauthenticated, status.Convert(GRPCError(err)).Code())

	err = xerrors.Errorf("ErrHTTP503: %w", ErrHTTP503)
	require.Equal(t, codes.Unavailable, status.Convert(GRPCError(err)).Code())

	err = xerrors.Errorf("ErrGRPCInternal: %w", ErrGRPCInternal)
	require.Equal(t, codes.Internal, status.Convert(GRPCError(err)).Code())

	err = xerrors.Errorf("ErrGRPCInvalidArgument: %w", ErrGRPCInvalidArgument)
	require.Equal(t, codes.InvalidArgument, status.Convert(GRPCError(err)).Code())

	err = xerrors.Errorf("ErrGRPCInvalidArgument, repo.ErrExists: %w", multierror.Append(ErrGRPCInvalidArgument, repo.ErrExists))
	require.Equal(t, codes.AlreadyExists, status.Convert(GRPCError(err)).Code())

	err = xerrors.Errorf("ErrHTTP502, repo.ErrExists: %w", multierror.Append(ErrHTTP502, repo.ErrExists))
	require.Equal(t, codes.AlreadyExists, status.Convert(GRPCError(err)).Code())

	err = xerrors.Errorf("xerrors.New(\"multi\"), sql.ErrNoRows): %w", multierror.Append(xerrors.New("multi"), sql.ErrNoRows))
	require.Equal(t, codes.NotFound, status.Convert(GRPCError(err)).Code())
}
