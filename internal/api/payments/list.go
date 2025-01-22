package payments

import (
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getPaymentList(ctx *gin.Context) {
	a := api.Get(ctx)
	
	payment, err := a.Queries.ListPaymentsOrderByCreated(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
