package processor_controller

import (
	"context"
	processor_domain "repo-stat/processor/internal/domain"
	processorpb "repo-stat/proto/processor"

	"google.golang.org/protobuf/types/known/emptypb"
)

type ProcessorService interface {
	GetRepoInfo(context.Context, string, string) ([]*processor_domain.RepoInfo, error)
	Ping(context.Context) (*processor_domain.Ping, error)
}

type processorController struct {
	ps ProcessorService
	processorpb.UnimplementedProcessorServer
}

func NewProcessorService(procserv ProcessorService) *processorController {
	return &processorController{
		ps: procserv,
	}
}

func (pc *processorController) GetInfoRepo(ctx context.Context, req *processorpb.GetInfoRepoRequest) (*processorpb.GetInfoRepoResponse, error) {
	resp, err := pc.ps.GetRepoInfo(ctx, req.GetRepo(), req.GetOwner())

	if err != nil {
		return nil, HandleErrorFromDomainToGRPC(err)
	}

	return &processorpb.GetInfoRepoResponse{
		Repoinfo: &processorpb.RepoInfo{
			Fullname:    resp[0].FullName,
			Description: resp[0].Description,
			Forks:       resp[0].Forks,
			Stargazers:  resp[0].Stargazers,
			Createdat:   resp[0].CreatedAt,
		},
	}, nil
}

func (pc *processorController) GetInfoRepositories(ctx context.Context, req *emptypb.Empty) (*processorpb.GetInfoRepositoriesResponse, error) {
	repos, err := pc.ps.GetRepoInfo(ctx, "", "")

	if err != nil {

		return nil, HandleErrorFromDomainToGRPC(err)
	}

	resp := make([]*processorpb.RepoInfo, len(repos))

	for i := range repos {
		resp[i] = &processorpb.RepoInfo{
			Fullname:    repos[i].FullName,
			Description: repos[i].Description,
			Forks:       repos[i].Forks,
			Stargazers:  repos[i].Stargazers,
			Createdat:   repos[i].CreatedAt,
		}
	}

	return &processorpb.GetInfoRepositoriesResponse{
		Repositoriesinfo: resp,
	}, nil
}

func (pc *processorController) Ping(ctx context.Context, req *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	responce, _ := pc.ps.Ping(ctx)
	return &processorpb.PingResponse{
		Reply: responce.Reply,
	}, nil
}
