package rpc

import (
	"database/sql"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"

	"github.com/eiicon-company/go-core/util/repo"
)

type (
	// Mux uses grpc mux
	Mux struct {
		GrpcMux      *grpc.Server
		GrpcEndpoint string
		GwMux        *runtime.ServeMux
		GwOpts       []grpc.DialOption
		GwEndpoint   string
	}

	// Muxer has gRPC configuration
	Muxer struct {
		API     *Mux
		Company *Mux
	}
)

// AsList simply list converter
func (m *Muxer) AsList() []*Mux {
	return []*Mux{m.API, m.Company}
}

// GeneralTimeout context timeout seconds
var GeneralTimeout = time.Second * 30

// GRPCError returns grpc status error
func GRPCError(err error) error {
	if _, ok := err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		return err
	}
	if errors.Cause(err) == sql.ErrNoRows {
		return status.Error(codes.NotFound, err.Error())
	}
	if errors.Cause(err) == repo.ErrExists {
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
