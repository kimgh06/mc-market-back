package user

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

func updateUser(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.WithJSON(err.Error()))
		return
	}

	if !permissions.RequireUserPermission(ctx, user, permissions.ManageUsers) {
		return
	}

	body := UpdateUserBody{}
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

	updated, err := a.Queries.UpdateUser(ctx, schema.UpdateUserParams{
		ID:          int64(id),
		Nickname:    nullable.PointerToString(body.Nickname),
		Permissions: nullable.PointerToInt32(body.Permissions),
		Cash:        nullable.PointerToInt32(body.Cash),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, responses.FullUserFromSchema(updated))
	return
}
