package products

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/godruoyi/go-snowflake"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strings"
)

func purchaseProduct(ctx *gin.Context) {
	a := api.Get(ctx)
	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	user := middlewares.GetUser(ctx)

	product, err := a.Queries.GetProductById(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.ProductNotFound.MakeJSON())
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		}
		return
	}

	currentPrice := GetProductPrice(product.Product.PriceDiscount, product.Product.Price)

	if user.Cash < currentPrice {
		ctx.AbortWithStatusJSON(http.StatusPaymentRequired, perrors.InsufficientFunds.MakeJSON(fmt.Sprintf("%d required", currentPrice)))
		return
	}

	err = a.Transaction(ctx, func(tx *sql.Tx, queries *schema.Queries) error {
		_, err := queries.CreatePurchase(ctx, schema.CreatePurchaseParams{
			ID:        int64(snowflake.ID()),
			Purchaser: user.ID,
			Product:   int64(id),
			Cost:      currentPrice,
		})
		if err != nil {
			return err
		}

		_, err = queries.UpdateUser(ctx, schema.UpdateUserParams{
			ID:   user.ID,
			Cash: sql.NullInt32{Int32: user.Cash - currentPrice, Valid: true},
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			ctx.AbortWithStatusJSON(http.StatusConflict, perrors.DuplicatePurchase.MakeJSON(err.Error()))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}

func getPurchase(ctx *gin.Context) {
	a := api.Get(ctx)
	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	user := middlewares.GetUser(ctx)

	_, err = a.Queries.GetPurchase(ctx, schema.GetPurchaseParams{
		Purchaser: user.ID,
		Product:   int64(id),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusOK, false)
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, true)
}
