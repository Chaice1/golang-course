package collectordomain

import (
	"context"
	"errors"

	collectorrespmodel "repo-stat/collector/internal/dto"
)

type RepoInfo struct {
	FullName    string
	Description string
	Stargazers  uint64
	Forks       uint64
	CreatedAt   string
}

var (
	ErrNotFound      = errors.New("NOT_FOUND")
	ErrInternalError = errors.New("INTERNAL_ERROR")
	ErrBadRequest    = errors.New("BAD_REQUEST")
)

type GitHubClient interface {
	GetRepoInfo(context.Context, string, string) (*collectorrespmodel.RepoInfo, error)
}
