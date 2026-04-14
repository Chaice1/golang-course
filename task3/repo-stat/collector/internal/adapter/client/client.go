package collectorclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	collectordomain "repo-stat/collector/internal/domain"

	collectorrespmodel "repo-stat/collector/internal/dto"
)

type GitHubApiClient struct{}

func (ghac *GitHubApiClient) GetRepoInfo(ctx context.Context, owner string, repo string) (*collectorrespmodel.RepoInfo, error) {

	client := http.Client{}
	url := "https://api.github.com/repos/" + owner + "/" + repo

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("request error:%w", collectordomain.ErrBadRequest)
	}

	req.Header.Set("User-Agent", "my-github-cli-tool")

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("github api call: %w", collectordomain.ErrInternalError)
	}
	defer func() { _ = resp.Body.Close() }()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, collectordomain.ErrNotFound
	case http.StatusInternalServerError:
		return nil, collectordomain.ErrInternalError
	}
	RepoInfoSlice, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read respbody: %w", collectordomain.ErrInternalError)
	}

	var RepInfo collectorrespmodel.RepoInfo
	err = json.Unmarshal(RepoInfoSlice, &RepInfo)
	log.Println(RepInfo)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse resp body: %w", collectordomain.ErrInternalError)
	}
	return &RepInfo, nil
}
