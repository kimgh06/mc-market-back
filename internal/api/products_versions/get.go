package products_versions

import (
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getProductVersion struct {
	ID         int32  `json:"id"`
	ProductID  int64  `json:"product_id"`
	VersionName string `json:"version_name"`
	Link       string `json:"link"`
	Index      int    `json:"index"`
	UpdatedAt  string `json:"updated_at"`
}

func getOneVersion(ctx *gin.Context) {
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

	log, err := a.Queries.GetOneProductVersion(ctx, int32(pid))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	response := getProductVersion{
		ID:         log.ID,
		ProductID:  log.ProductID,
		VersionName: log.VersionName,
		Link:       log.Link,
		Index:      log.Index,
		UpdatedAt:  log.UpdatedAt.String(),
	}

	ctx.JSON(http.StatusOK, response)
}