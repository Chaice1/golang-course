package processor_controller

import (
	"context"
	"log/slog"
	processor_controller_errors "repo-stat/processor/internal/controller/errors"
	processor_domain "repo-stat/processor/internal/domain"
	processorpb "repo-stat/proto/processor"
	"time"
)

type ProcessorService interface {
	GetRepoInfo(context.Context, []*processor_domain.GetRepoInfoRequestBody) ([]*processor_domain.RepoInfo, error)
	Ping(context.Context) string
}

type processorController struct {
	ps ProcessorService
	processorpb.UnimplementedProcessorServer
	log *slog.Logger
}

func NewProcessorService(procserv ProcessorService, log *slog.Logger) *processorController {
	return &processorController{
		ps:  procserv,
		log: log,
	}
}

func (pc *processorController) GetInfoRepo(ctx context.Context, req *processorpb.GetInfoRepoRequest) (*processorpb.GetInfoRepoResponse, error) {
	resp, err := pc.ps.GetRepoInfo(ctx, []*processor_domain.GetRepoInfoRequestBody{&processor_domain.GetRepoInfoRequestBody{Repo: req.GetRepo(), Owner: req.GetOwner()}})

	if err != nil {
		return nil, processor_controller_errors.HandleErrorFromDomainToGRPC(err, pc.log, "GetInfoRepo")
	}

	return &processorpb.GetInfoRepoResponse{
		Repoinfo: &processorpb.RepoInfo{
			Fullname:    resp[0].FullName,
			Description: resp[0].Description,
			Forks:       resp[0].Forks,
			Stargazers:  resp[0].Stargazers,
			Createdat:   resp[0].CreatedAt.String(),
		},
	}, nil
}

func (pc *processorController) GetInfoRepositories(ctx context.Context, req *processorpb.GetInfoRepositoriesRequest) (*processorpb.GetInfoRepositoriesResponse, error) {

	GetRepoInfoSlice := make([]*processor_domain.GetRepoInfoRequestBody, len(req.GetReqs()))

	req_slice := req.GetReqs()
	for i := range req_slice {
		GetRepoInfoSlice[i] = &processor_domain.GetRepoInfoRequestBody{
			Owner: req_slice[i].GetOwner(),
			Repo:  req_slice[i].GetRepo(),
		}
	}
	repos, err := pc.ps.GetRepoInfo(ctx, GetRepoInfoSlice)

	if err != nil {
		return nil, processor_controller_errors.HandleErrorFromDomainToGRPC(err, pc.log, "GetInfoRepositories")
	}

	resp := make([]*processorpb.RepoInfo, len(repos))

	for i := range repos {
		resp[i] = &processorpb.RepoInfo{
			Fullname:    repos[i].FullName,
			Description: repos[i].Description,
			Forks:       repos[i].Forks,
			Stargazers:  repos[i].Stargazers,
			Status:      repos[i].Status,
			Createdat:   repos[i].CreatedAt.Format(time.RFC3339),
		}
	}

	return &processorpb.GetInfoRepositoriesResponse{
		Repositoriesinfo: resp,
	}, nil
}

func (pc *processorController) Ping(ctx context.Context, req *processorpb.PingRequest) (*processorpb.PingResponse, error) {
	responce := pc.ps.Ping(ctx)
	return &processorpb.PingResponse{
		Reply: responce,
	}, nil
}
