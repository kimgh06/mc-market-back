package products

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/godruoyi/go-snowflake"
	"html/template"
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

	sanitizedDetails := buffer.String()

	product, err := a.Queries.CreateProduct(ctx, schema.CreateProductParams{
		ID:          int64(snowflake.ID()),
		Creator:     int64(body.Creator),
		Name:        body.Name,
		Description: body.Description,
		Usage:       body.Usage,
		Details:     sanitizedDetails,
		Category:    body.Category,
		Price:       body.Price,
		Tags:        body.Tags,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, responses.ProductFromSchema(&schema.GetProductByIdRow{
		Product: *product,
		User:    *user,
		Count:   0,
	}))
	return
}
