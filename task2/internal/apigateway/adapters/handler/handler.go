package apigatewayhandler

import (
	"context"
	"net/http"

	apigatewayrepoinfomodel "github.com/Chaice1/golang-course/task2/internal/apigateway/adapters/repoinfomodel"
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

func (agh *apigatewayHandler) GetRepoInfo(c *gin.Context) {
	var req apigatewayrepoinfomodel.GetRepoInfoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": apigatewaydomain.BadRequest,
		})
		return
	}

	resp, err := agh.uag.GetInfoRep(c.Request.Context(), req.Owner, req.Repo)

	var httpstatus int

	switch err {
	case apigatewaydomain.BadRequest:
		httpstatus = http.StatusBadRequest
	case apigatewaydomain.NotFound:
		httpstatus = http.StatusNotFound
	case apigatewaydomain.InternalError:
		httpstatus = http.StatusInternalServerError
	}

	if err != nil {
		c.JSON(httpstatus, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "repo is found successfully",
		"repoinfo": apigatewayrepoinfomodel.GetRepoInfoResp{
			FullName:    resp.FullName,
			Description: resp.Description,
			Stargazers:  resp.Stargazers,
			Forks:       resp.Forks,
			CreatedAt:   resp.CreatedAt,
		},
	})

}
