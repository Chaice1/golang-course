package subsciber_repository

import (
	"context"
	generated_db "repo-stat/subscriber/internal/adapter/repository/gen"
	subscriber_domain "repo-stat/subscriber/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	q *generated_db.Queries
}

func NewRepo(p *pgxpool.Pool) *repo {
	return &repo{
		q: generated_db.New(p),
	}
}

func (r *repo) CreateSubscription(ctx context.Context, repo string, owner string) error {
	err := r.q.CreateSubscription(ctx, generated_db.CreateSubscriptionParams{
		RepoName:  repo,
		OwnerName: owner,
	})

	if err != nil {
		return subscriber_domain.ErrInternalError
	}
	return nil
}

func (r *repo) DeleteSubscription(ctx context.Context, repo string, owner string) error {
	err := r.q.DeleteSubscription(ctx, generated_db.DeleteSubscriptionParams{
		RepoName:  repo,
		OwnerName: owner,
	})

	if err != nil {
		return subscriber_domain.ErrInternalError
	}
	return nil
}

func (r *repo) GetSubscriptions(ctx context.Context) ([]*subscriber_domain.Subscription, error) {
	responce, err := r.q.GetSubscriptions(ctx)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, subscriber_domain.ErrNotFound
		default:
			return nil, subscriber_domain.ErrInternalError
		}
	}

	subscriptions := make([]*subscriber_domain.Subscription, len(responce))

	for i := range responce {
		subscriptions[i] = &subscriber_domain.Subscription{
			Id:        responce[i].ID,
			RepoName:  responce[i].RepoName,
			OwnerName: responce[i].OwnerName,
			CreatedAt: responce[i].CreatedAt.Time.String(),
		}
	}

	return subscriptions, nil
}

func (r *repo) GetSubscription(ctx context.Context, repo string, owner string) (*subscriber_domain.Subscription, error) {
	responce, err := r.q.GetSubscription(ctx, generated_db.GetSubscriptionParams{
		RepoName:  repo,
		OwnerName: owner,
	})

	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, subscriber_domain.ErrNotFound
		default:
			return nil, subscriber_domain.ErrInternalError
		}
	}
	return &subscriber_domain.Subscription{
		Id:        responce.ID,
		OwnerName: responce.OwnerName,
		RepoName:  responce.RepoName,
		CreatedAt: responce.CreatedAt.Time.String(),
	}, nil
}
