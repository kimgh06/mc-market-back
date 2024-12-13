package products

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/pkg/permissions"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func uploadImage(ctx *gin.Context) {
	a := api.Get(ctx)
	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	product, err := a.Queries.GetProductById(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.ProductNotFound.MakeJSON())
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		}
		return
	}

	user := middlewares.GetUser(ctx)
	if !permissions.CheckUserPermission(permissions.UserPermission(user.Permissions), permissions.ManageProducts) && product.Product.Creator != user.ID {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON())
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedFormFile.MakeJSON(err.Error()))
		return
	}

	err = ctx.SaveUploadedFile(file, filepath.Join(a.Config.Storage.ImagesPath, "products", strconv.FormatUint(id, 10)))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}

func getImage(ctx *gin.Context) {
	a := api.Get(ctx)
	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	path := filepath.Join(a.Config.Storage.ImagesPath, "products", strconv.FormatUint(id, 10))

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		ctx.Status(http.StatusNoContent)
		return
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	contentType := http.DetectContentType(bytes)

	ctx.Header("Content-Type", contentType)
	ctx.Status(http.StatusOK)
	_, _ = ctx.Writer.Write(bytes)
}
