package banner

import (
	"errors"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateBannerImage struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	IndexNum int    `json:"index_num"`
}

type CreateBanner struct {
	Title    string `json:"title"`
	LinkURL	string `json:"link_url"`
	IndexNum int    `json:"index_num"`
}

func uploadImage(ctx *gin.Context) {
	a := api.Get(ctx)

	body := CreateBanner{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.FailedFormFile.MakeJSON(err.Error()))
		return
	}

	imagePath := time.Now().Format("20060102150405")

	err = ctx.SaveUploadedFile(file, filepath.Join(a.Config.Storage.ImagesPath, "banners", imagePath))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
		return
	}

	// create banner in database
	_, err = a.Queries.CreateBanner(ctx, schema.CreateBannerParams{ 
		Title:    body.Title,
		ImageURL: "/v1/banner/image/" + imagePath,
		LinkURL: body.LinkURL,
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

	image, err := a.Queries.GetBanner(ctx, int64(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if image == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.BannerNotFound.MakeJSON("Banner not found"))
		return
	}

	body := UpdateBannerImage{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	body.ImageURL = image.ImageURL

	newFile, err := ctx.FormFile("file") 
	if err == nil {
		imagePath := time.Now().Format("20060102150405")
		err = ctx.SaveUploadedFile(newFile, filepath.Join(a.Config.Storage.ImagesPath, "banners", imagePath))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedStorage.MakeJSON(err.Error()))
			return
		}
		body.ImageURL = imagePath
	}

	if body.Title == "" {
		body.Title = image.Title
	}

	if body.IndexNum == 0 {
		body.IndexNum = image.IndexNum
	}

	if body.LinkURL == "" {
		body.LinkURL = image.LinkURL
	}

	_, err = a.Queries.UpdateBanner(ctx, schema.UpdateBannerParams{
		ID:       int64(id),
		Title:    body.Title,
		ImageURL: body.ImageURL,
		LinkURL:  body.LinkURL,
		IndexNum: body.IndexNum,
	})
}

func getImage(ctx *gin.Context) {
	a := api.Get(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	image, err := a.Queries.GetBanner(ctx, int64(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if image == nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.BannerNotFound.MakeJSON("Banner not found"))
		return
	}

	ctx.JSON(http.StatusOK, schema.Banner{
		ID:       image.ID,
		Title:    image.Title,
		ImageURL: image.ImageURL,
		LinkURL:  image.LinkURL,
		IndexNum: image.IndexNum,
	})
}

func getImageFromUrl(ctx *gin.Context) {
	path := ctx.Param("path")
	println(path)
	
	if path == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.BannerNotFound.MakeJSON("Invalid path"))
		return
	}

	imagePath := filepath.Join(api.Get(ctx).Config.Storage.ImagesPath, "banners", path)
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

	banners, err := a.Queries.ListBanners(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	var response []schema.Banner
	for _, banner := range banners {
		response = append(response, schema.Banner{
			ID:       banner.ID,
			Title:    banner.Title,
			ImageURL: banner.ImageURL,
			LinkURL:  banner.LinkURL,
			IndexNum: banner.IndexNum,
			CreatedAt: banner.CreatedAt,
		})
	}

	if response == nil {
		response = []schema.Banner{}
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

	err = a.Queries.DeleteBanner(ctx, int64(id))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}