package user

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/nullable"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/internal/utilities"
	"net/http"
	"strconv"
	"time"
)

type unclaimedRevenuesElement struct {
	Cost            int       `json:"cost"`
	ProductID       string    `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductPrice    int       `json:"product_price"`
	ProductDiscount *int32    `json:"product_discount"`
	PurchasedAt     time.Time `json:"date"`
}

func getUnclaimedRevenues(ctx *gin.Context) {
	a := api.Get(ctx)
	userID := middlewares.GetUserID(ctx)

	result, err := a.Queries.GetUnclaimedPurchasesOfUser(ctx, int64(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusOK, []int{})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	mapped := utilities.Map(result, func(t *schema.GetUnclaimedPurchasesOfUserRow) unclaimedRevenuesElement {
		return unclaimedRevenuesElement{
			Cost:            int(t.Cost),
			ProductID:       strconv.FormatUint(uint64(t.Product.ID), 10),
			ProductName:     t.Product.Name,
			ProductPrice:    int(t.Product.Price),
			ProductDiscount: nullable.Int32ToPointer(t.Product.PriceDiscount),
			PurchasedAt:     t.PurchasedAt,
		}
	})

	ctx.JSON(http.StatusOK, mapped)
}
