package payments

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"
)

func approvePayment(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	orderId := api.GetUUIDFromParams(ctx, "orderId")
	if orderId == uuid.Nil {
		return
	}

	amountString := ctx.Query("amount")
	amount, err := strconv.Atoi(amountString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidQuery.MakeJSON(err.Error()))
		return
	}

	payment, err := a.Queries.GetPaymentByOrderId(ctx, orderId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.PaymentNotFound.MakeJSON(err.Error()))
		return
	}

	if payment.Agent != user.ID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON())
		return
	}

	if payment.Approved {
		ctx.AbortWithStatus(http.StatusOK)
		return
	}

	if int(payment.Amount) != amount {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.Mismatching.MakeJSON())
	}

	if err = a.Transaction(ctx, func(tx *sql.Tx, queries *schema.Queries) error {
		_, err = queries.ApprovePayment(ctx, payment.ID)
		if err != nil {
			return err
		}

		user, err = queries.UpdateUser(ctx, schema.UpdateUserParams{
			ID:   user.ID,
			Cash: sql.NullInt32{Int32: user.Cash + int32(amount), Valid: true},
		})
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}
