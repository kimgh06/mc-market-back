package article_likes

import (
	"database/sql"
	"errors"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

func deleteArticleLike(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	id, err := api.GetUint64FromParam(ctx, "article_id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	_, err = a.Queries.GetArticleLikeById(ctx, schema.GetArticleLikeParams{ArticleID: uint64(id), UserID: uint64(user.ID)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.ArticleLikeNotFound.MakeJSON(err.Error()))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if err = a.Queries.DeleteArticleLike(ctx, schema.DeleteArticleLikeParams{ArticleID: uint64(id), UserID: uint64(user.ID)}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}
