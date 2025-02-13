package products_update_log

import (
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getProductUpdateLogListResponse struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UpdateAt  string `json:"update_at"`
}

func getUpdateList(ctx *gin.Context) {
	a := api.Get(ctx)
	productID := ctx.Param("product_id")
	if productID == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("Invalid product_id parameter"))
		return
	}

	pid, err := strconv.ParseInt(productID, 10, 64)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("product_id must be an integer"))
		return
	}

	logs, err := a.Queries.ListUpdateLogsByProductID(ctx, pid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	var response []getProductUpdateLogListResponse
	for _, log := range logs {
		response = append(response, getProductUpdateLogListResponse{
			ID:        int64(log.ID),
			ProductID: int64(log.ProductID),
			Title:     log.Title,
			Content:   log.Content,
			UpdateAt:  log.UpdatedAt.String(),
		})
	}

	if response == nil {
		response = []getProductUpdateLogListResponse{}
	}

	ctx.JSON(http.StatusOK, response)
}