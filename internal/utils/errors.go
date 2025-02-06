package er

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// APIError represents an error with an associated HTTP status code
type APIError struct {
	StatusCode int
	Message    string
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}

// my custom implementation
func NewError(statusCode int, message string) error {
	return MapErrorToGRPCStatus(&APIError{
		StatusCode: statusCode,
		Message:    message,
	})
}

func MapErrorToGRPCStatus(err error) error {
	if apiErr, ok := err.(*APIError); ok {
		grpcCode := MapHTTPToGRPCCode(apiErr.StatusCode)
		return status.Errorf(codes.Code(grpcCode), "%s", apiErr.Message)
	}
	return err
}

func MapHTTPToGRPCCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusInternalServerError:
		return codes.Internal
	case http.StatusServiceUnavailable:
		return codes.Unavailable
	default:
		return codes.Unknown
	}
}
