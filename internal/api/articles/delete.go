package articles

import (
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"
)

func deleteArticle(ctx *gin.Context) {
	a := api.Get(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	if err = a.Queries.DeleteArticle(ctx, int64(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}
