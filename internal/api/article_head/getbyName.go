package article_head

import (
	"database/sql"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getArticleHeadbyName struct {
	Name string `json:"name"`
}

func getHeadByName(ctx *gin.Context) {
	a := api.Get(ctx)

	name, boool := ctx.Params.Get("name")
	if !boool {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("name is required"))
		return
	}

	id, err := a.Queries.GetArticleHeadByName(ctx, schema.GetArticleHeadByNameParams{Name: name})
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.ArticleHeadNotFound.MakeJSON("article head not found"))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, id)
}