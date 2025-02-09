package products_update_log

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateProductsUpdateLogBody struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Title     string `json:"title" binding:"required,max=100"`
	Content   string `json:"content" binding:"required"`
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

	body := UpdateProductsUpdateLogBody{}
	if err = ctx.ShouldBind(&body); err != nil {
		// failed to bind body, abort
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	// Assume you have a corresponding query generated.

	_, err = a.Queries.UpdateUpdateLog(ctx, schema.UpdateProductUpdateLogParams{
		ID:        int32(id),
		ProductID: body.ProductID,
		Title:     body.Title,
		Content:   body.Content,
		UpdatedAt: time.Now(),
	})
	
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}
	
	ctx.Status(http.StatusNoContent)
}