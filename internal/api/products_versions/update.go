package products_versions

import (
	"fmt"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateProductVersion struct {
	VersionName string `json:"version_name" binding:"required"`
	Link       string `json:"link" binding:"required"`
	Index      int    `json:"index"`
}

func update(ctx *gin.Context) {
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

	body := UpdateProductVersion{}
	if err = ctx.ShouldBindBodyWithJSON(&body); err != nil {
		// failed to bind body, abort
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	fmt.Println(body)

	// Assume you have a corresponding query generated.	
	_, err = a.Queries.UpdateProductVersion(ctx, schema.UpdateProductVersionParams{
		ID:         int32(id),
		VersionName: body.VersionName,
		Link:       body.Link,
		Index:      body.Index,
	})	
	
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}
	
	ctx.Status(http.StatusNoContent)
}