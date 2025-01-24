package payments

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	paymentKey := ctx.Query("paymentKey")
	if paymentKey == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidQuery.MakeJSON("paymentKey is required"))
		return
	}

	paymentSecret := os.Getenv("TOSS_PAYMENTS_SECRET_KEY")
	if paymentSecret == "" {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.MissingEnvironmentVariable.MakeJSON("TOSS_PAYMENTS_SECRET_KEY is required"))
		return
	}
	encrypted := "Basic " + base64.StdEncoding.EncodeToString([]byte(paymentSecret+":"))

	// request to payment service
	url := "https://api.tosspayments.com/v1/payments/confirm"
	payload := map[string]interface{}{
		"orderId":    orderId,
		"amount":     amount,
		"paymentKey": paymentKey,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedToMarshal.MakeJSON(err.Error()))
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedToCreateRequest.MakeJSON(err.Error()))
		return
	}
	req.Header.Set("Authorization", encrypted)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedToSendRequest.MakeJSON(err.Error()))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		ctx.AbortWithStatusJSON(resp.StatusCode, perrors.PaymentFailed.MakeJSON(string(bodyBytes)))
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
