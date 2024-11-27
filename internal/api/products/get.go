package products

import (
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/api/responses"
	"maple/internal/perrors"
	"net/http"
	"strconv"
	"strings"
)

func getProduct(ctx *gin.Context) {
	a := api.Get(ctx)

	idString := strings.Trim(ctx.Param("id"), "")
	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON("invalid snowflake", idString))
		return
	}

	product, err := a.Queries.GetProductById(ctx, int64(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, responses.ProductFromSchema(product))
}
