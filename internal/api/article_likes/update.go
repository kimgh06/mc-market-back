package article_likes

import (
	"database/sql"
	"errors"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateArticleLike struct {
	Like *bool `json:"like" binding:"required"` // 수정: exists -> required
}

func updateArticleLike(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	var body UpdateArticleLike
	if err := ctx.ShouldBindJSON(&body); err != nil { // 수정: ShouldBind -> ShouldBindJSON
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	articleIDStr := ctx.Param("article_id")
	if articleIDStr == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("article_id is required"))
		return
	}

	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("article_id must be an integer"))
		return
	}

	_, err = a.Queries.GetArticleLikeById(ctx, schema.GetArticleLikeParams{ArticleID: uint64(articleID), UserID: uint64(user.ID)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.ArticleLikeNotFound.MakeJSON("like/dislike not found for this article"))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if body.Like == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("The 'like' field is required and must be a boolean"))
		return
	}

	if err = a.Queries.UpdateArticleLike(ctx, schema.UpdateArticleLikeParams{
		ArticleID: uint64(articleID),
		UserID:    uint64(user.ID),
		Kind:      *body.Like,
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}
