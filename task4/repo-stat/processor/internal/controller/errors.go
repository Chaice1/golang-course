package processor_controller

import (
	processor_domain "repo-stat/processor/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleErrorFromDomainToGRPC(err error) error {
	switch err {
	case processor_domain.ErrBadRequest:
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	case processor_domain.ErrInternalError:
		return status.Error(codes.Internal, codes.Internal.String())
	case processor_domain.ErrNotFound:
		return status.Error(codes.NotFound, codes.NotFound.String())
	default:
		return status.Error(codes.Unknown, codes.Unknown.String())
	}
}
