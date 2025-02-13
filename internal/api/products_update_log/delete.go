package products_update_log

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func deleteLog(ctx *gin.Context) {
	a := api.Get(ctx)
	// check user permission
	user := middlewares.GetUser(ctx)
	if user.Permissions != 2147483647 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to delete a log"))
		return
	}

	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("Invalid id parameter"))
		return
	}

	// Assume you have a corresponding query generated.
	err = a.Queries.DeleteUpdateLog(ctx, int32(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusNoContent)
}
