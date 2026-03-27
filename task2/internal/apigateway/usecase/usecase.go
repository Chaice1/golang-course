package apigatewayusecase

import (
	"context"

	apigatewaydomain "github.com/Chaice1/golang-course/task2/internal/apigateway/domain"
)

type usecaseApiGateway struct {
	cc apigatewaydomain.CollectorClient
}

func NewUsecaseApiGateway(cc apigatewaydomain.CollectorClient) *usecaseApiGateway {
	return &usecaseApiGateway{
		cc: cc,
	}
}

func (uag *usecaseApiGateway) GetInfoRep(ctx context.Context, owner string, repo string) (*apigatewaydomain.RepoInfo, error) {
	return uag.cc.GetRepoInfo(ctx, owner, repo)
}
