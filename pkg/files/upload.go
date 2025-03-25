package files

import (
	"bytes"
	"encoding/json"
	"io"
	"maple/internal/perrors"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func UploadAndReturnURL(ctx *gin.Context, file *multipart.FileHeader) string {
	url := os.Getenv("IMG_URL")

	if url == "" {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Failed to get image url"))
		return ""
	}

	// create multipart form
	formBody := &bytes.Buffer{}
	writer := multipart.NewWriter(formBody)
	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return ""
	}

	fileContent, err := file.Open()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return ""
	}
	defer fileContent.Close()

	_, err = io.Copy(part, fileContent)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return ""
	}
	writer.Close()

	// api request to imgbb
	req, _ := http.NewRequest("POST", url+"/upload", formBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return ""
	}
	defer res.Body.Close()

	// get response from imgbb
	var img_response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&img_response)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return ""
	}

	imagePath, ok := img_response["url"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Failed to get image url"))
		return ""
	}

	return url + imagePath
}

func ReadImage(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}