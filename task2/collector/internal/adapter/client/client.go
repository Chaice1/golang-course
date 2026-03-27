package collectorclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	collectordomain "github.com/Chaice1/golang-course/task2/collector/internal/domain"

	collectorrespmodel "github.com/Chaice1/golang-course/task2/collector/internal/adapter/resp_model"
)

type GitHubApiClient struct{}

func (ghac *GitHubApiClient) GetRepoInfo(ctx context.Context, owner string, repo string) (*collectorrespmodel.RepoInfo, error) {

	client := http.Client{}
	url := "https://api.github.com/repos/" + owner + "/" + repo

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("request error:%w", collectordomain.BadRequest)
	}

	req.Header.Set("User-Agent", "my-github-cli-tool")

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("github api call: %w", collectordomain.InternalError)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return nil, collectordomain.ErrorNotFound
	case http.StatusInternalServerError:
		return nil, collectordomain.InternalError
	}

	RepoInfoSlice, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldn't read respbody: %w", collectordomain.InternalError)
	}

	var RepInfo collectorrespmodel.RepoInfo
	err = json.Unmarshal(RepoInfoSlice, &RepInfo)

	if err != nil {
		return nil, fmt.Errorf("couldn't parse resp body: %w", collectordomain.InternalError)
	}
	return &RepInfo, nil
}
