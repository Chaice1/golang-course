package collectorhandler

import (
	"context"

	collectordomain "repo-stat/collector/internal/domain"

	collectorpb "repo-stat/proto/collector"

	"google.golang.org/protobuf/types/known/emptypb"
)

type UsecaseCollectorService interface {
	GetInfoRepo(context.Context, string, string) ([]*collectordomain.RepoInfo, error)
}

type colletorHandler struct {
	collectorpb.UnimplementedCollectorServer
	ucs UsecaseCollectorService
}

func NewHandler(ucs UsecaseCollectorService) *colletorHandler {
	return &colletorHandler{
		ucs: ucs,
	}
}

func (h *colletorHandler) GetInfoRepo(ctx context.Context, req *collectorpb.GetInfoRepoRequest) (*collectorpb.GetInfoRepoResponse, error) {

	RepoInfo, err := h.ucs.GetInfoRepo(ctx, req.GetOwner(), req.GetRepo())

	if err != nil {
		return nil, HandleErrorsFromDomainToGRPC(err)
	}

	return &collectorpb.GetInfoRepoResponse{
		Repoinfo: &collectorpb.RepoInfo{
			Fullname:    RepoInfo[0].FullName,
			Description: RepoInfo[0].Description,
			Stargazers:  RepoInfo[0].Stargazers,
			Forks:       RepoInfo[0].Forks,
			Createdat:   RepoInfo[0].CreatedAt,
		},
	}, nil
}

func (h *colletorHandler) GetInfoRepositories(ctx context.Context, req *emptypb.Empty) (*collectorpb.GetInfoRepositoriesResponse, error) {
	RepositoriesInfo, err := h.ucs.GetInfoRepo(ctx, "", "")

	if err != nil {
		return nil, HandleErrorsFromDomainToGRPC(err)
	}

	resp := make([]*collectorpb.RepoInfo, len(RepositoriesInfo))

	for i := range RepositoriesInfo {
		resp[i] = &collectorpb.RepoInfo{
			Fullname:    RepositoriesInfo[i].FullName,
			Description: RepositoriesInfo[i].Description,
			Forks:       RepositoriesInfo[i].Forks,
			Stargazers:  RepositoriesInfo[i].Stargazers,
			Createdat:   RepositoriesInfo[i].CreatedAt,
		}
	}

	return &collectorpb.GetInfoRepositoriesResponse{
		Repositoriesinfo: resp,
	}, nil

}
