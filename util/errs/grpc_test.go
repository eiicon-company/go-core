package errs

import (
	"database/sql"
	"testing"

	"golang.org/x/xerrors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/go-multierror"

	"github.com/eiicon-company/go-core/util/repo"
)

func TestGRPCError(t *testing.T) {
	t.Helper()

	var err error

	err = xerrors.Errorf("ErrHTTPUnauthorized: %w", ErrHTTPUnauthorized)
	if status.Convert(GRPCError(err)).Code() != codes.Unauthenticated {
		t.Errorf("fatal get error code: %+#v", err)
	}

	err = xerrors.Errorf("ErrHTTP503: %w", ErrHTTP503)
	if status.Convert(GRPCError(err)).Code() != codes.Unavailable {
		t.Errorf("fatal get error code: %+#v", err)
	}

	err = xerrors.Errorf("ErrGRPCInternal: %w", ErrGRPCInternal)
	if status.Convert(GRPCError(err)).Code() != codes.Internal {
		t.Errorf("fatal get error code: %+#v", err)
	}

	err = xerrors.Errorf("ErrGRPCInvalidArgument: %w", ErrGRPCInvalidArgument)
	if status.Convert(GRPCError(err)).Code() != codes.InvalidArgument {
		t.Errorf("fatal get error code: %+#v", err)
	}

	err = xerrors.Errorf("ErrGRPCInvalidArgument, repo.ErrExists: %w", multierror.Append(ErrGRPCInvalidArgument, repo.ErrExists))
	if status.Convert(GRPCError(err)).Code() != codes.AlreadyExists {
		t.Errorf("fatal get error code: %+#v", err)
	}

	err = xerrors.Errorf("ErrHTTP502, repo.ErrExists: %w", multierror.Append(ErrHTTP502, repo.ErrExists))
	if status.Convert(GRPCError(err)).Code() != codes.AlreadyExists {
		t.Errorf("fatal get error code: %+#v", err)
	}

	err = xerrors.Errorf("xerrors.New(\"multi\"), sql.ErrNoRows): %w", multierror.Append(xerrors.New("multi"), sql.ErrNoRows))
	if status.Convert(GRPCError(err)).Code() != codes.NotFound {
		t.Errorf("fatal get error code: %+#v", err)
	}
}
