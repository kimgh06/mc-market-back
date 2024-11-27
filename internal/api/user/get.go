package user

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"
)

func getUser(ctx *gin.Context) {
	a := api.Get(ctx)
	targetId, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		return
	}

	//user := middlewares.GetUser(ctx)

	target, err := a.Queries.GetUserById(ctx, int64(targetId))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.UserNotFound.MakeJSON())
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		}
		return
	}

	var res userGetResponse
	res.ID = uint64(target.ID)
	if target.Nickname.Valid {
		res.Nickname = &target.Nickname.String
	}

	ctx.JSON(http.StatusOK, res)
}
