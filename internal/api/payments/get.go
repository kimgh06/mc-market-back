package payments

import (
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"net/http"
)

func getPayment(ctx *gin.Context) {
	a := api.Get(ctx)

	orderId := api.GetUUIDFromParams(ctx, "orderId")

	payment, err := a.Queries.GetPaymentByOrderId(ctx, orderId)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, payment)
}
