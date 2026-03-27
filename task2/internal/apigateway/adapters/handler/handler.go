package apigatewayhandler

import (
	"context"
	"net/http"

	apigatewaymodel "github.com/Chaice1/golang-course/task2/internal/apigateway/adapters/model"
	apigatewaydomain "github.com/Chaice1/golang-course/task2/internal/apigateway/domain"
	"github.com/gin-gonic/gin"
)

type UsecaseApiGateway interface {
	GetInfoRep(context.Context, string, string) (*apigatewaydomain.RepoInfo, error)
}
type apigatewayHandler struct {
	uag UsecaseApiGateway
}

func NewApiGatewayHandler(uag UsecaseApiGateway) *apigatewayHandler {
	return &apigatewayHandler{
		uag: uag,
	}
}

// @Summary      Get Repository Info
// @Description  Returns information about a  GitHub repository with the help of Collector
// @Tags         Gateway
// @Accept       json
// @Produce      json
// @Param        owner  path string  true  "GitHub Owner Name"
// @Param        repo   path string  true  "Repository Name"
// @Success      200 {object} apigatewaymodel.SuccessResponce "success"
// @Failure      404 {string} string "NOT_FOUND"
// @Failure      500 {string} string "INTERNAL_ERROR"
// @Failure      400 {string} string "BAD_REQUEST"
// @Router       /get_repo_info/{owner}/{repo} [get]
func (agh *apigatewayHandler) GetRepoInfo(c *gin.Context) {

	owner := c.Param("owner")
	repo := c.Param("repo")
	if len(owner) == 0 || len(repo) == 0 {
		c.JSON(http.StatusBadRequest, apigatewaymodel.ErrorResponce{
			Message: apigatewaydomain.BadRequest.Error(),
		})
		return
	}

	resp, err := agh.uag.GetInfoRep(c.Request.Context(), owner, repo)

	httpstatus := http.StatusInternalServerError

	switch err {
	case apigatewaydomain.BadRequest:
		httpstatus = http.StatusBadRequest
	case apigatewaydomain.NotFound:
		httpstatus = http.StatusNotFound
	}

	if err != nil {
		c.JSON(httpstatus, apigatewaymodel.ErrorResponce{
			Message: err.Error(),
		})
		return
	}

	responce := apigatewaymodel.GetRepoInfoResp{
		FullName:    resp.FullName,
		Description: resp.Description,
		Stargazers:  resp.Stargazers,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
	}

	c.JSON(http.StatusOK, apigatewaymodel.SuccessResponce{
		Message: "success",
		RepInfo: responce,
	})

}
