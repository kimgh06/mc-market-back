package products_update_log

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateProductUpdateLog struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func createProductUpdateLog(ctx *gin.Context) {
	a := api.Get(ctx)
	// Check user permission
	user := middlewares.GetUser(ctx)
	if user.Permissions != 2147483647 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to create product update log"))
		return
	}

	body := CreateProductUpdateLog{}
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	productIDStr := ctx.Param("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("Invalid product_id"))
		return
	}

	_, err = a.Queries.CreateUpdateLog(ctx, schema.CreateProductUpdateLogParams{
		ProductID: int64(productID),
		Title:     body.Title,
		Content:   body.Content,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}
