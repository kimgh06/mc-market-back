package user

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

func listUsers(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	if !permissions.RequireUserPermission(ctx, user, permissions.ListUsers) {
		return
	}

	offset := utilities.Clamp(api.QueryIntDefault(ctx, "offset", 0), 0, math.MaxInt)
	limit := utilities.Clamp(api.QueryIntDefault(ctx, "limit", 20), 0, 20)

	users, err := a.Queries.ListUsers(ctx, schema.ListUsersParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	fullUsers := utilities.Map(users, func(u *schema.User) responses.FullUser {
		return responses.FullUserFromSchema(u)
	})

	ctx.JSON(http.StatusOK, fullUsers)
}
