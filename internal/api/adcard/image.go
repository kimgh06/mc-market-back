package adcard

import (
	"errors"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/pkg/files"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateAdcardImage struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	IndexNum int    `json:"index_num"`
}

type CreateAdcard struct {
	Title    string `json:"title"`
	LinkURL  string `json:"link_url"`
	IndexNum int    `json:"index_num"`
}

func uploadImage(ctx *gin.Context) {
	a := api.Get(ctx)

	body := CreateAdcard{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

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

	// create adcard in database
	_, err = a.Queries.CreateAdcard(ctx, schema.CreateAdcardParams{
		Title:    body.Title,
		ImageURL: returnedURL,
		LinkURL:  body.LinkURL,
		IndexNum: body.IndexNum,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}

func updateImage(ctx *gin.Context) {
	a := api.Get(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	adcard, err := a.Queries.GetAdcard(ctx, int64(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if adcard == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.BannerNotFound.MakeJSON("Adcard not found"))
		return
	}

	body := UpdateAdcardImage{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	body.ImageURL = adcard.ImageURL

	newFile, err := ctx.FormFile("file")
	if err == nil {
		imagePath := time.Now().Format("20060102150405")
		err = ctx.SaveUploadedFile(newFile, filepath.Join(a.Config.Storage.ImagesPath, "adcards", imagePath))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
			return
		}
		body.ImageURL = imagePath
	}

	if newFile != nil {
		returnedURL := files.UploadAndReturnURL(ctx, newFile)
		if returnedURL == "" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Failed to upload image"))
			return
		}
		body.ImageURL = returnedURL
	}

	if body.Title == "" {
		body.Title = adcard.Title
	}

	if body.IndexNum == 0 {
		body.IndexNum = adcard.IndexNum
	}

	if body.LinkURL == "" {
		body.LinkURL = adcard.LinkURL
	}

	_, err = a.Queries.UpdateAdcard(ctx, schema.UpdateAdcardParams{
		ID:       int64(id),
		Title:    body.Title,
		ImageURL: body.ImageURL,
		LinkURL:  body.LinkURL,
		IndexNum: body.IndexNum,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
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

	adcard, err := a.Queries.GetAdcard(ctx, int64(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if adcard == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.BannerNotFound.MakeJSON("Adcard not found"))
		return
	}

	ctx.JSON(http.StatusOK, schema.Adcard{
		ID:       adcard.ID,
		Title:    adcard.Title,
		ImageURL: adcard.ImageURL,
		LinkURL:  adcard.LinkURL,
		IndexNum: adcard.IndexNum,
	})
}

func getImageFromUrl(ctx *gin.Context) {
	path := ctx.Param("path")

	if path == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.BannerNotFound.MakeJSON("Invalid path"))
		return
	}

	imagePath := filepath.Join(api.Get(ctx).Config.Storage.ImagesPath, "adcards", path)
	if _, err := os.Stat(imagePath); errors.Is(err, os.ErrNotExist) {
		ctx.Status(http.StatusNoContent)
		return
	}

	bytes, err := os.ReadFile(imagePath)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	contentType := http.DetectContentType(bytes)

	ctx.Header("Content-Type", contentType)
	ctx.Status(http.StatusOK)
	_, _ = ctx.Writer.Write(bytes)
}

func getListImage(ctx *gin.Context) {
	a := api.Get(ctx)

	adcards, err := a.Queries.ListAdcards(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	var response []schema.Adcard
	for _, ad := range adcards {
		response = append(response, schema.Adcard{
			ID:        ad.ID,
			Title:     ad.Title,
			ImageURL:  ad.ImageURL,
			LinkURL:   ad.LinkURL,
			IndexNum:  ad.IndexNum,
			CreatedAt: ad.CreatedAt,
		})
	}

	if response == nil {
		response = []schema.Adcard{}
	}

	ctx.JSON(http.StatusOK, response)
}

func deleteImage(ctx *gin.Context) {
	a := api.Get(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	err = a.Queries.DeleteAdcard(ctx, int64(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}