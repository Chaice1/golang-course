package usecase

import (
	"context"
	"errors"
	subscriber_domain "repo-stat/subscriber/internal/domain"
)

type SubscriberRepo interface {
	GetSubscriptions(context.Context) ([]*subscriber_domain.Subscription, error)
	CreateSubscription(context.Context, string, string) error
	DeleteSubscription(context.Context, string, string) error
	GetSubscription(context.Context, string, string) (*subscriber_domain.Subscription, error)
}

type CollectorClient interface {
	GetRepoInfo(context.Context, string, string) error
}

type subscriberUsecase struct {
	sr SubscriberRepo
	cc CollectorClient
}

func NewSubscriberUsecase(sr SubscriberRepo, cc CollectorClient) *subscriberUsecase {
	return &subscriberUsecase{
		sr: sr,
		cc: cc,
	}
}

func (su *subscriberUsecase) GetSubscriptions(ctx context.Context) ([]*subscriber_domain.Subscription, error) {
	return su.sr.GetSubscriptions(ctx)
}

func (su *subscriberUsecase) DeleteSubscription(ctx context.Context, repo string, owner string) error {
	return su.sr.DeleteSubscription(ctx, repo, owner)
}

func (su *subscriberUsecase) CreateSubscription(ctx context.Context, repo string, owner string) error {

	_, err := su.sr.GetSubscription(ctx, repo, owner)

	if err == nil {
		return subscriber_domain.ErrSubscriptionAlreadyExists
	}

	if !errors.Is(err, subscriber_domain.ErrNotFound) {
		return err
	}

	err = su.cc.GetRepoInfo(ctx, owner, repo)
	if err != nil {
		return err
	}

	return su.sr.CreateSubscription(ctx, repo, owner)

}
