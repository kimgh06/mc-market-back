package products

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"maple/internal/api"
	"maple/internal/api/responses"
	"maple/internal/middlewares"
	"maple/internal/nullable"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/pkg/permissions"
	"net/http"
)

func updateProduct(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.WithJSON(err.Error()))
		return
	}

	if !permissions.RequireUserPermission(ctx, user, permissions.ManageProducts) {
		return
	}

	body := UpdateProductBody{}
	if err = ctx.ShouldBind(&body); err != nil {
		// failed to bind body, abort
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedValidate.MakeJSON(err.Error()))
		return
	}

	product, err := a.Queries.UpdateProduct(ctx, schema.UpdateProductParams{
		ID:          int64(id),
		Creator:     nullable.UPointerToInt64(body.Creator),
		Name:        nullable.PointerToString(body.Name),
		Description: nullable.PointerToString(body.Description),
		Usage:       nullable.PointerToString(body.Usage),
		Category:    nullable.PointerToString(body.Category),
		Price:       nullable.PointerToInt32(body.Price),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, responses.ProductFromSchema(product))
	return
}
