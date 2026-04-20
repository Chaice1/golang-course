package processor_adapter

import (
	processor_domain "repo-stat/processor/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorHandleFromGRPCToDomain(err error) error {
	status, ok := status.FromError(err)
	if !ok {
		return processor_domain.ErrInternalError
	}

	switch status.Code() {
	case codes.InvalidArgument:
		return processor_domain.ErrBadRequest
	case codes.NotFound:
		return processor_domain.ErrNotFound
	default:
		return processor_domain.ErrInternalError
	}
}
