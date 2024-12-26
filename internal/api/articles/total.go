package articles

import (
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"
)

func countArticles(ctx *gin.Context) {
	a := api.Get(ctx)

	count, err := a.Queries.CountArticles(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, count)
}
