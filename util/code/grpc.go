package code

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

// GRPC2HTTP returns HTTP Status
func GRPC2HTTP(code interface{}) int {
	switch v := code.(type) {
	default:
		return 500
	case codes.Code:
		return grpc2http(v)
	case int:
		return GRPC2HTTP(codes.Code(v))
	}
}

func grpc2http(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Aborted:
		return http.StatusConflict
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Canceled:
		return 499 // Be returned this code by the HTTP Server such as Nginx.
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	default:
		return http.StatusInternalServerError // TODO: Connectivity issues
	}
}
