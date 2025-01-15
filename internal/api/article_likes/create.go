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

type CreateArticleLike struct {
	Like *bool `json:"like" binding:"required"` // 'required'는 필드가 존재하고 값이 있어야 유효
}

func createArticleLike(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	var body CreateArticleLike
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	if body.Like == nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("like is required"))
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
		if !errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
			return
		}
	}

	// find the existing row
	_, err = a.Queries.GetArticleLikeById(ctx, schema.GetArticleLikeParams{ArticleID: uint64(articleID), UserID: uint64(user.ID)})
	// not exist
	if errors.Is(err, sql.ErrNoRows) {
		// Create a new row
		if err = a.Queries.CreateArticleLike(ctx, schema.CreateArticleLikeParams{
			ArticleID: uint64(articleID),
			UserID:    uint64(user.ID),
			Kind:      *body.Like,
		}); err != nil {
			// Update the existing row
			if err = a.Queries.UpdateArticleLike(ctx, schema.UpdateArticleLikeParams{
				ArticleID: uint64(articleID),
				UserID:    uint64(user.ID),
				Kind:      *body.Like,
			}); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
				return
			}
		}
	} else{
		if err = a.Queries.DeleteArticleLike(ctx, schema.DeleteArticleLikeParams{ArticleID: uint64(articleID), UserID: uint64(user.ID)}); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
			return
		}
	}

	ctx.Status(http.StatusOK)
}