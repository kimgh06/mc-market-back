package products

import (
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"html/template"
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	product, err := a.Queries.GetProductById(ctx, int64(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if !permissions.CheckUserPermission(permissions.UserPermission(user.Permissions), permissions.ManageProducts) && product.Product.Creator != user.ID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON())
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

	tmpl, err := template.New("product.details").Parse(body.Details)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.UnknownInternalError.MakeJSON(err.Error()))
		return
	}
	buffer := new(bytes.Buffer)

	if err = tmpl.Execute(buffer, body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidHTML.MakeJSON(err.Error()))
		return
	}

	if !permissions.CheckUserPermission(permissions.UserPermission(user.Permissions), permissions.ManageProducts) {
		casted := uint64(user.ID)
		body.Creator = &casted
	}

	updated, err := a.Queries.UpdateProduct(ctx, schema.UpdateProductParams{
		ID:          int64(id),
		Details:     sql.NullString{String: buffer.String(), Valid: true},
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

	ctx.JSON(http.StatusOK, responses.ProductFromSchema(&schema.GetProductByIdRow{
		Product: *updated,
		User:    *user,
		Count:   0,
	}))
	return
}
