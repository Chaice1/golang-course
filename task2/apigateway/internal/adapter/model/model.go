package apigatewaymodel

type GetRepoInfoResp struct {
	FullName    string `json:"full_name"`
	Description string `json:"descritpion"`
	Stargazers  uint64 `json:"stargazers"`
	Forks       uint64 `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type ErrorResponce struct {
	Message string `json:"message"`
}

type SuccessResponce struct {
	Message string          `json:"message"`
	RepInfo GetRepoInfoResp `json:"repo_info"`
}
