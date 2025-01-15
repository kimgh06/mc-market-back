package article_head

import (
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getArticleHeadListResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func getHeadList(ctx *gin.Context) {
	a := api.Get(ctx)

	heads, err := a.Queries.GetArticleHeadList(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	var response []getArticleHeadListResponse
	for _, head := range heads {
		response = append(response, getArticleHeadListResponse{
			ID:   int64(head.ID),
			Name: head.Name,
		})
	}
	
	if response == nil {
		response = []getArticleHeadListResponse{}
	}

	ctx.JSON(http.StatusOK, response)
}