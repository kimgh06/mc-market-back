package article_likes

import (
	"database/sql"
	"errors"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getArticleLikesResponse struct {
	Dislikes int64 `json:"dislikes"`
	Likes    int64 `json:"likes"`
}

func getArticleLikesAndDisLikes(ctx *gin.Context) {
	a := api.Get(ctx)

	id, err := api.GetUint64FromParam(ctx, "article_id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	likes, dislikes, err := a.Queries.GetArticleLikesAndDisLikesCount(ctx, schema.GetArticleLikesCountParams{ArticleID: uint(id)})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.ArticleNotFound.MakeJSON(err.Error()))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, getArticleLikesResponse{
		Dislikes: dislikes,
		Likes:    likes,
	})

}
