package collectorusecase

import (
	"context"
	"log/slog"

	collectordomain "repo-stat/collector/internal/domain"

	"golang.org/x/sync/errgroup"
)

type collectorService struct {
	ghc collectordomain.GitHubClient
	sc  collectordomain.SubscriberClient
}

func NewCollectorService(ghc collectordomain.GitHubClient, sc collectordomain.SubscriberClient) *collectorService {
	return &collectorService{
		ghc: ghc,
		sc:  sc,
	}
}

func (cs *collectorService) GetInfoRepo(ctx context.Context, owner string, repo string) ([]*collectordomain.RepoInfo, error) {

	if owner != "" && repo != "" {

		RepoInfo, err := cs.ghc.GetRepoInfo(ctx, owner, repo)
		if err != nil {
			return nil, err
		}

		return []*collectordomain.RepoInfo{RepoInfo}, nil
	}

	subscriptions, err := cs.sc.GetSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	semaphore := make(chan struct{}, 10)

	g, ctxx := errgroup.WithContext(ctx)

	RepositoriesInfo := make([]*collectordomain.RepoInfo, len(subscriptions))

	for i := range subscriptions {
		idx := i
		semaphore <- struct{}{}
		g.Go(func() error {
			defer func() { <-semaphore }()

			RepoInfo, err := cs.ghc.GetRepoInfo(ctxx, subscriptions[idx].OwnerName, subscriptions[idx].RepoName)
			if err != nil {
				return err
			}
			RepositoriesInfo[idx] = RepoInfo
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		slog.Error("message", "error", err)
		return nil, err
	}

	return RepositoriesInfo, nil
}
