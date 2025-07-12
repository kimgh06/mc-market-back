package settlements

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"regexp"
)

var monthRegex = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])$`)

func getSettlementByUserIdAndMonth(ctx *gin.Context) {
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

	logrus.Debugf("getSettlementByUserIdAndMonth: user_id=%d, month=%s, settlement=%+v", user.ID, month, settlement)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, settlement)
}

func getSettlementsByUserIdAndMonths(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)
	start_month := ctx.Query("start_month")
	end_month := ctx.Query("end_month")

	if start_month == "" || end_month == "" {
		ctx.JSON(400, gin.H{"error": "start_month and end_month are required"})
		return
	}
	if !monthRegex.MatchString(start_month) || !monthRegex.MatchString(end_month) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "month must be in 'YYYY-MM' format"})
		return
	}

	settlements, err := a.Queries.GetSettlementsByUserIdAndMonths(ctx, schema.GetSettlementsParams{
		UserId:     user.ID,
		StartMonth: start_month,
		EndMonth:   end_month,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, settlements)
}
