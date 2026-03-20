package apigatewayclient

import (
	"context"

	collectorpb "github.com/Chaice1/golang-course/task2/gen"
	apigatewaydomain "github.com/Chaice1/golang-course/task2/internal/apigateway/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type collectorClient struct {
	collectorpb.CollectorClient
}

func NewCollectorClient(cc collectorpb.CollectorClient) *collectorClient {
	return &collectorClient{CollectorClient: cc}
}

func (cc *collectorClient) GetRepoInfo(ctx context.Context, owner string, repo string) (*apigatewaydomain.RepoInfo, error) {

	repoinfo, err := cc.GetInfoRepo(ctx, &collectorpb.GetInfoRepoRequest{
		Owner: owner,
		Repo:  repo,
	})

	status, ok := status.FromError(err)
	if !ok {
		return nil, apigatewaydomain.InternalError
	}

	switch status.Code() {
	case codes.InvalidArgument:
		return nil, apigatewaydomain.BadRequest
	case codes.NotFound:
		return nil, apigatewaydomain.NotFound
	case codes.Internal:
		return nil, apigatewaydomain.InternalError
	}

	return &apigatewaydomain.RepoInfo{
		FullName:    repoinfo.Fullname,
		Description: repoinfo.Description,
		Stargazers:  repoinfo.Stargazers,
		Forks:       repoinfo.Forks,
		CreatedAt:   repoinfo.Createdat,
	}, nil
}
