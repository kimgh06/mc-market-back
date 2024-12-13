package user

import (
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/middlewares"
	"net/http"
)

type unclaimedRevenuesResponse struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

func getUnclaimedRevenues(ctx *gin.Context) {
	a := api.Get(ctx)
	userID := middlewares.GetUserID(ctx)

	total, err := a.Queries.GetUnclaimedRevenuesOfUser(ctx, int64(userID))
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, unclaimedRevenuesResponse{
		Count: int(total.Count),
		Total: int(total.Coalesce.(int64)),
	})
}
