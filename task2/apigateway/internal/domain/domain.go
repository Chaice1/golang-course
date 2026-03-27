package apigatewaydomain

import (
	"context"
	"errors"
)

type RepoInfo struct {
	FullName    string
	Description string
	Stargazers  uint64
	Forks       uint64
	CreatedAt   string
}

type CollectorClient interface {
	GetRepoInfo(context.Context, string, string) (*RepoInfo, error)
}

var (
	NotFound      = errors.New("NOT_FOUND")
	InternalError = errors.New("INTERNAL_ERROR")
	BadRequest    = errors.New("BAD_REQUEST")
)
