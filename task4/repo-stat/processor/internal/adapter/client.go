package processor_adapter

import (
	"context"
	processor_domain "repo-stat/processor/internal/domain"
	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type collectorClient struct {
	cc   collectorpb.CollectorClient
	conn *grpc.ClientConn
}

func NewCollectorClient(address string) (*collectorClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, processor_domain.ErrInternalError
	}

	return &collectorClient{
		cc:   collectorpb.NewCollectorClient(conn),
		conn: conn,
	}, nil

}

func (cc *collectorClient) GetRepoInfo(ctx context.Context, repo string, owner string) ([]*processor_domain.RepoInfo, error) {

	if repo != "" && owner != "" {
		resp, err := cc.cc.GetInfoRepo(ctx, &collectorpb.GetInfoRepoRequest{
			Owner: owner,
			Repo:  repo,
		})

		if err != nil {
			return nil, ErrorHandleFromGRPCToDomain(err)
		}

		RepoInfo := resp.GetRepoinfo()
		return []*processor_domain.RepoInfo{
			&processor_domain.RepoInfo{
				FullName:    RepoInfo.Fullname,
				Description: RepoInfo.Description,
				Forks:       RepoInfo.Forks,
				Stargazers:  RepoInfo.Stargazers,
				CreatedAt:   RepoInfo.Createdat,
			},
		}, nil
	}

	resp, err := cc.cc.GetInfoRepositories(ctx, &emptypb.Empty{})

	if err != nil {
		return nil, ErrorHandleFromGRPCToDomain(err)
	}

	repos := make([]*processor_domain.RepoInfo, len(resp.GetRepositoriesinfo()))

	for i, item := range resp.GetRepositoriesinfo() {

		repos[i] = &processor_domain.RepoInfo{
			FullName:    item.Fullname,
			Description: item.Description,
			Forks:       item.Forks,
			Stargazers:  item.Stargazers,
			CreatedAt:   item.Createdat,
		}
	}

	return repos, nil
}

func (cc *collectorClient) Ping(ctx context.Context) (*processor_domain.Ping, error) {
	return &processor_domain.Ping{Reply: "Pong"}, nil
}
