package products_versions

import (
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getVersionList(ctx *gin.Context) {
	a := api.Get(ctx)
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("Invalid product_id parameter"))
		return
	}

	pid, err := strconv.ParseInt(productID, 10, 64)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("product_id must be an integer"))
		return
	}

	logs, err := a.Queries.ListProductVersionsByProductID(ctx, pid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	var response []getProductVersion
	for _, log := range logs {
		response = append(response, getProductVersion{
			ID:         log.ID,
			ProductID:  log.ProductID,
			VersionName: log.VersionName,
			Link:       log.Link,
			Index:      log.Index,
			UpdatedAt:  log.UpdatedAt.String(),
		})
	}

	ctx.JSON(http.StatusOK, response)
}