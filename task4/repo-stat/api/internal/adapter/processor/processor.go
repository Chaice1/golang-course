package processor

import (
	"context"
	"log/slog"
	adapter_errors "repo-stat/api/internal/adapter/errors"
	"repo-stat/api/internal/domain"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type processorClient struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pc   processorpb.ProcessorClient
}

func NewProcessorClient(addres string, log *slog.Logger) (*processorClient, error) {

	conn, err := grpc.NewClient(addres, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	client := processorpb.NewProcessorClient(conn)

	return &processorClient{
		log:  log,
		conn: conn,
		pc:   client,
	}, nil
}
func (pc *processorClient) Ping(ctx context.Context) domain.PingStatus {
	_, err := pc.pc.Ping(ctx, &processorpb.PingRequest{})
	if err != nil {
		pc.log.Error("processor ping failed", "error", err)
		return domain.PingStatusDown
	}
	return domain.PingStatusUp

}
func (pc *processorClient) GetInfoRepo(ctx context.Context, owner string, repo string) ([]*domain.RepoInfo, error) {
	if owner != "" && repo != "" {
		resp, err := pc.pc.GetInfoRepo(ctx, &processorpb.GetInfoRepoRequest{
			Owner: owner,
			Repo:  repo,
		})

		if err != nil {
			return nil, adapter_errors.ErrorHandleFromGRPCToDomainWithLog(err, pc.log, "GetInfoRepo")
		}

		repoInfo := resp.GetRepoinfo()
		return []*domain.RepoInfo{&domain.RepoInfo{
			FullName:    repoInfo.Fullname,
			Description: repoInfo.Description,
			Forks:       repoInfo.Forks,
			Stargazers:  repoInfo.Stargazers,
			CreatedAt:   repoInfo.Createdat,
		}}, nil
	}

	resp, err := pc.pc.GetInfoRepositories(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, adapter_errors.ErrorHandleFromGRPCToDomainWithLog(err, pc.log, "GetInfoRepo")
	}

	repositoriesInfo := make([]*domain.RepoInfo, len(resp.GetRepositoriesinfo()))

	for i, item := range resp.GetRepositoriesinfo() {
		repositoriesInfo[i] = &domain.RepoInfo{
			FullName:    item.Fullname,
			Description: item.Description,
			Forks:       item.Forks,
			Stargazers:  item.Stargazers,
			CreatedAt:   item.Createdat,
		}
	}
	return repositoriesInfo, nil
}

func (pc *processorClient) mapRepInfoToDomain(repoInfo *processorpb.RepoInfo) *domain.RepoInfo {
	return &domain.RepoInfo{
		FullName:    repoInfo.Fullname,
		Description: repoInfo.Description,
		Forks:       repoInfo.Forks,
		Stargazers:  repoInfo.Stargazers,
		CreatedAt:   repoInfo.Createdat,
	}
}

func (pc *processorClient) Close() error {
	return pc.conn.Close()
}
