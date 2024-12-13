package products

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/pkg/permissions"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func uploadFile(ctx *gin.Context) {
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

	directoryPath := filepath.Join(a.Config.Storage.ContentsPath, "products", strconv.FormatUint(id, 10))

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedFormFile.MakeJSON(err.Error()))
		return
	}

	err = ctx.SaveUploadedFile(file, filepath.Join(directoryPath, file.Filename))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	existingFiles, err := os.ReadDir(directoryPath)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	for _, existingFile := range existingFiles {
		if existingFile.Name() != file.Filename {
			_ = os.Remove(filepath.Join(directoryPath, existingFile.Name()))
		}
	}

	ctx.Status(http.StatusOK)
}

func getFile(ctx *gin.Context) {
	a := api.Get(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	user := middlewares.GetUser(ctx)

	if !permissions.CheckUserPermission(permissions.UserPermission(user.Permissions), permissions.ManageProducts) {
		_, err = a.Queries.GetPurchase(ctx, schema.GetPurchaseParams{
			Purchaser: user.ID,
			Product:   int64(id),
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				ctx.AbortWithStatusJSON(http.StatusPaymentRequired, perrors.PurchaseNotFound.MakeJSON(err.Error()))
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
			}
			return
		}
	}

	path := filepath.Join(a.Config.Storage.ContentsPath, "products", strconv.FormatUint(id, 10))

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		ctx.Status(http.StatusNoContent)
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil || len(entries) < 1 {
		ctx.Status(http.StatusNoContent)
		return
	}

	file := entries[0]
	path = filepath.Join(path, file.Name())

	stat, err := os.Stat(path)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=\""+file.Name()+"\"")
	ctx.Header("Content-Length", strconv.FormatInt(stat.Size(), 10))
	ctx.Header("Access-Control-Allow-Headers", "*")
	ctx.Header("Access-Control-Expose-Headers", "*")
	ctx.Status(http.StatusOK)
	_, _ = ctx.Writer.Write(bytes)
}
