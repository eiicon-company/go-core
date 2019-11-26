package rpc

import (
	"database/sql"
	"time"

	"golang.org/x/xerrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/eiicon-company/go-core/util/repo"
)

// GeneralTimeout context timeout seconds
var GeneralTimeout = time.Second * 30

// GRPCError returns grpc status error
func GRPCError(err error) error {
	if _, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return err
	}
	if xerrors.Is(err, sql.ErrNoRows) {
		return status.Error(codes.NotFound, err.Error())
	}
	if xerrors.Is(err, repo.ErrExists) {
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return status.Error(codes.Unknown, err.Error())
	}

	return nil
}
