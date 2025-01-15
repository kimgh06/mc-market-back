package article_head

import (
	"database/sql"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getArticleHeadbyId struct {
	Id int64 `json:"id"`
}

func getHeadById(ctx *gin.Context) {
	a := api.Get(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	head, err := a.Queries.GetArticleHeadByID(ctx, schema.GetArticleHeadByIDParams{ID: int(id)})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.ArticleHeadNotFound.MakeJSON("article head not found"))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, head)
}