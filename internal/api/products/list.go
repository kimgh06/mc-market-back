package products

import (
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/api/responses"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/internal/utilities"
	"maple/pkg/permissions"
	"math"
	"net/http"
)

func listProducts(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	if !permissions.RequireUserPermission(ctx, user, permissions.ManageProducts) {
		return
	}

	offset := utilities.Clamp(api.QueryIntDefault(ctx, "offset", 0), 0, math.MaxInt)
	limit := utilities.Clamp(api.QueryIntDefault(ctx, "limit", 20), 0, 20)

	products, err := a.Queries.ListProducts(ctx, schema.ListProductsParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	fullProducts := utilities.Map(products, func(u *schema.ListProductsRow) responses.ProductWithShortUser {
		return responses.ProductWithShortUserFromSchema(u)
	})

	ctx.JSON(http.StatusOK, fullProducts)
}
