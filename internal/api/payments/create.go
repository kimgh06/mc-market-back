package payments

import (
	"github.com/gin-gonic/gin"
	"github.com/godruoyi/go-snowflake"
	"github.com/google/uuid"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
)

type CreatePayment struct {
	OrderID uuid.UUID `json:"order_id"`
	Amount  int       `json:"amount"`
}

func createPayment(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	body := CreatePayment{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	_, err := a.Queries.CreatePayment(ctx, schema.CreatePaymentParams{
		ID:      int64(snowflake.ID()),
		Agent:   user.ID,
		OrderID: body.OrderID,
		Amount:  int32(body.Amount),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}
