package articles

import (
	"github.com/gin-gonic/gin"
	"github.com/godruoyi/go-snowflake"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
)

type CreateArticle struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func createArticle(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	body := CreateArticle{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	_, err := a.Queries.CreateArticle(ctx, schema.CreateArticleParams{
		ID:      int64(snowflake.ID()),
		Title:   body.Title,
		Content: body.Content,
		Author:  user.ID,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}
