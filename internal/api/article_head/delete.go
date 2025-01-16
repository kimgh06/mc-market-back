package article_head

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func deleteHead(ctx *gin.Context) {
	a := api.Get(ctx)
	// check user permission
	user := middlewares.GetUser(ctx)
	if user.Permissions != 2147483647 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to delete article head"))
		return
	}

	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("Invalid id parameter"))
		return
	}

	err = a.Queries.DeleteArticleHead(ctx, schema.DeleteArticleHeadParams{
		ID: id,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusNoContent)
}
