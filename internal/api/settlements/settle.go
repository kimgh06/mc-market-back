package settlements

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/schema"
	"net/http"
)

func settleByUserIdAndMonth(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)
	month := ctx.Query("month")

	if month == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "month is required"})
		return
	}
	if !monthRegex.MatchString(month) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "month must be in 'YYYY-MM' format"})
		return
	}

	settlement, err := a.Queries.GetSettlementByUserIdAndMonth(ctx, schema.GetSettlementParams{
		UserId: user.ID,
		Month:  month,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if settlement.Settled != false {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "settlement already exists for this month"})
		return
	}
	err = a.Transaction(ctx, func(tx *sql.Tx, queries *schema.Queries) error {
		_, err = queries.Settle(ctx, schema.SettleParams{
			UserId: user.ID,
			Month:  month,
		})
		if err != nil {
			return err
		}
		return nil
	})

	ctx.JSON(http.StatusOK, settlement)
}
