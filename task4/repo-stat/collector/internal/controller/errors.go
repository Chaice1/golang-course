package collectorhandler

import (
	"errors"
	collectordomain "repo-stat/collector/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleErrorsFromDomainToGRPC(err error) error {

	switch {
	case errors.Is(err, collectordomain.ErrBadRequest):
		return status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	case errors.Is(err, collectordomain.ErrNotFound):
		return status.Error(codes.NotFound, codes.NotFound.String())
	default:
		return status.Error(codes.Internal, codes.Internal.String())
	}
}
