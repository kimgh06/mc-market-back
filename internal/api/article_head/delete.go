package article_head

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteArticleHead struct {
	Id 	int `json:"id"`
}

func deleteHead(ctx *gin.Context) {
	a := api.Get(ctx)
	// check user permission
	user := middlewares.GetUser(ctx)
	if user.Permissions != 2147483647 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to create article head"))
		return
	}
	
	body := deleteArticleHead{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	err := a.Queries.DeleteArticleHead(ctx, schema.DeleteArticleHeadParams{
		ID:   body.Id,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}