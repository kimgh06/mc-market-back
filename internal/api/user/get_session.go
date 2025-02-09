package user

import (
	"database/sql"
	"errors"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getSessionUser(ctx *gin.Context) {
	a := api.Get(ctx)

	user := middlewares.GetUser(ctx)
	user, err := a.Queries.GetUserById(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.UserNotFound.MakeJSON())
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		}
		return
	}

	var res userSessionGetResponse
	res.ID = strconv.Itoa(int(user.ID))
	if user.Nickname.Valid {
		res.Nickname = &user.Nickname.String
	}
	res.Permissions = user.Permissions
	res.Cash = uint(user.Cash)

	ctx.JSON(http.StatusOK, res)
}
