package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetRepoInfoHandler
// @Summary      Получить инфо о репозитории
// @Description  Запрашивает данные о звездах, форках и коммитах
// @Tags         repository
// @Produce      json
// @Param 		 url	 query     string  true	"Полная ссылка на репозиторий (например, https://github.com/golang/go)"
// @Success      200     {object}  dto.RepoInfo
// @Failure      400     {object}  string "bad request"
// @Failure		 404     {object}  string "repo not found"
// @Failure 	 500	 {object}  string "internal error"
// @Router       /api/repositories/info [get]
func NewGetInfoRepoHandler(log *slog.Logger, agu *usecase.ApiGatewayUsecase, eh *dto.ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url_str := r.URL.Query().Get("url")

		if url_str == "" {
			log.Error("require url", "error", errors.New("bad_request"))
			eh.CreateErrorResponce(log, w, http.StatusBadRequest, "require url")
			return
		}

		parsed_url, err := url.Parse(url_str)

		if err != nil {
			log.Error("url format is invalud", "error", err)
			eh.CreateErrorResponce(log, w, http.StatusBadRequest, "url format is invalid")
			return
		}

		path := strings.Trim(parsed_url.Path, "/")
		path_slice := strings.Split(path, "/")

		if len(path_slice) != 2 {
			log.Error("not enough arguments for request", "error", errors.New("bad_request"))
			eh.CreateErrorResponce(log, w, http.StatusBadRequest, "not enough arguments for request")
			return
		}

		repo_info, err := agu.GetInfoRep(r.Context(), path_slice[0], path_slice[1])
		if err != nil {
			s, ok := status.FromError(err)
			if !ok {
				log.Error("unknown grpc error", "error", errors.New("internal_error"))
				eh.CreateErrorResponce(log, w, http.StatusBadRequest, "unknown grpc error")
				return
			}
			switch s.Code() {
			case codes.InvalidArgument:
				log.Error("bad request", "error", codes.InvalidArgument)
				eh.CreateErrorResponce(log, w, http.StatusBadRequest, "bad request")
				return
			case codes.NotFound:
				log.Error("not found", "error", codes.NotFound)
				eh.CreateErrorResponce(log, w, http.StatusNotFound, "repo not found")
				return
			default:
				log.Error("internal error", "error", codes.Internal)
				eh.CreateErrorResponce(log, w, http.StatusInternalServerError, "internal error")
				return
			}
		}

		resp := dto.RepoInfo{
			FullName:    repo_info.FullName,
			Description: repo_info.Description,
			Forks:       repo_info.Forks,
			Stargazers:  repo_info.Stargazers,
			CreatedAt:   repo_info.CreatedAt,
		}

		w.Header().Add("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			log.Error("failed to write RepoInfo", "error", err)
		}
	}
}
