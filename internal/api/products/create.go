package products

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/godruoyi/go-snowflake"
	"maple/internal/api"
	"maple/internal/api/responses"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/pkg/permissions"
	"net/http"
)

func createProduct(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	body := CreateProductBody{}
	err := ctx.ShouldBind(&body)
	if err != nil {
		// failed to bind body, abort
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	if body.Creator == 0 {
		body.Creator = uint64(user.ID)
	}

	// If creator is specified, and it is not request user, requester must have ManageProduct permission in order to create product
	if body.Creator != uint64(user.ID) {
		if !permissions.RequireUserPermission(ctx, user, permissions.ManageProducts) {
			return
		}
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedValidate.MakeJSON(err.Error()))
		return
	}

	product, err := a.Queries.CreateProduct(ctx, schema.CreateProductParams{
		ID:          int64(snowflake.ID()),
		Creator:     int64(body.Creator),
		Name:        body.Name,
		Description: sql.NullString{String: body.Description, Valid: true},
		Usage:       sql.NullString{String: body.Usage, Valid: true},
		Category:    body.Category,
		Price:       body.Price,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, responses.ProductFromSchema(product))
	return
}
