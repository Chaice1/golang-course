package processor_domain

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

type Ping struct {
	Reply string
}

type GetRepoInfoRequestBody struct {
	Repo  string
	Owner string
}

var (
	ErrNotFound      = errors.New("NOT_FOUND")
	ErrInternalError = errors.New("INTERNAL_ERROR")
	ErrBadRequest    = errors.New("BAD_REQUEST")
)

type CollectorClient interface {
	GetRepoInfo(context.Context, string, string) ([]*RepoInfo, error)
	Ping(context.Context) (*Ping, error)
}
