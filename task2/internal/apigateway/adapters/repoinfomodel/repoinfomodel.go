package apigatewayrepoinfomodel

type GetRepoInfoReq struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type GetRepoInfoResp struct {
	FullName    string `json:"full_name"`
	Description string `json:"descritpion"`
	Stargazers  uint64 `json:"stargazers"`
	Forks       uint64 `json:"forks"`
	CreatedAt   string `json:"created_at"`
}
