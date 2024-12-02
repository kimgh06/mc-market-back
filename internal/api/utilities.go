package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"maple/internal/perrors"
	"net/http"
	"strconv"
)

func GetUUIDFromParams(ctx *gin.Context, key string) uuid.UUID {
	idString := ctx.Param(key)
	id, err := uuid.Parse(idString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON(err.Error()))
		return uuid.Nil
	}

	return id
}

func GetUint64FromParam(ctx *gin.Context, key string) (uint64, error) {
	valueString := ctx.Param(key)
	value, err := strconv.ParseUint(valueString, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON(err.Error()))
		return 0, err
	}

	return value, nil
}

func QueryIntDefault(ctx *gin.Context, key string, def int) int {
	value, err := strconv.Atoi(ctx.Query(key))
	if err != nil {
		return def
	}

	return value
}

func QueryStringDefault(ctx *gin.Context, key string, def string) string {
	v, exists := ctx.GetQuery(key)
	if !exists {
		return def
	}

	return v
}

func QueryUUIDDefault(ctx *gin.Context, key string, def uuid.UUID) uuid.UUID {
	parsed, err := uuid.Parse(ctx.Query(key))
	if err != nil {
		return def
	}

	return parsed
}
