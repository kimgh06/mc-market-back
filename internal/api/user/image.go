package user

import (
	"database/sql"
	"errors"
	"fmt"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/pkg/files"
	"net/http"

	"github.com/gin-gonic/gin"
)

func uploadImage(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedFormFile.MakeJSON(err.Error()))
		return
	}
	
	returnedURL := files.UploadAndReturnURL(ctx, file)
	if returnedURL == "" {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Failed to upload image"))
		return
	}

	err = a.Queries.UploadUserImage(ctx, schema.UploadUserImageParams{
		ID:       int64(user.ID),
		ImageURL: sql.NullString{String: returnedURL, Valid: true},
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	// err = ctx.SaveUploadedFile(file, filepath.Join(a.Config.Storage.ImagesPath, "avatars", strconv.FormatUint(uint64(user.ID), 10)))
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
	// 	return
	// }

	ctx.Status(http.StatusOK)
}

func getImage(ctx *gin.Context) {
	a := api.Get(ctx)
	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	imageURL, err := a.Queries.GetUserImage(ctx, int64(id))
	fmt.Println(imageURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.NotFound.MakeJSON())
		} else {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		}
		return
	}

	bytes, err := files.ReadImage(imageURL)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	contentType := http.DetectContentType(bytes)

	ctx.Header("Content-Type", contentType)
	ctx.Status(http.StatusOK)
	_, _ = ctx.Writer.Write(bytes)
}
