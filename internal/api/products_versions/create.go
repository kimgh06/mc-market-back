package products_versions

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateVersionParam struct {
	VersionName string `json:"version_name"`
	Link        string `json:"link"`
	Index       int    `json:"index"`
}

func createVersion(ctx *gin.Context) {
	a := api.Get(ctx)
	// Check user permission
	user := middlewares.GetUser(ctx)
	if user.Permissions != 2147483647 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to create product update log"))
		return
	}

	body := CreateVersionParam{}
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

	_, err = a.Queries.CreateProductVersion(ctx, schema.CreateProductVersionParams{
		ProductID:  int64(productID),
		VersionName: body.VersionName,
		Link:       body.Link,
		Index:      0,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}
