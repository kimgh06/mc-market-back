package products_update_log

import (
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getProductUpdateLogResponse struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UpdateAt  string `json:"update_at"`
}

func getUpdateOne(ctx *gin.Context) {
	a := api.Get(ctx)
	id:= ctx.Param("id")
	
	if id == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("Invalid id parameter"))
		return
	}

	pid, err := strconv.Atoi(id)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON("id must be an integer"))
		return
	}

	log, err := a.Queries.GetOneUpdateLog(ctx, int32(pid))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	response := getProductUpdateLogResponse{
		ID:        int64(log.ID),
		ProductID: int64(log.ProductID),
		Title:     log.Title,
		Content:   log.Content,
		UpdateAt:  log.UpdatedAt.String(),
	}

	ctx.JSON(http.StatusOK, response)
}