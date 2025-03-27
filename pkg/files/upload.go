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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return ""
	}

	var img_response map[string]interface{}
	if err = json.Unmarshal(body, &img_response); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Invalid response format: "+err.Error()))
		return ""
	}

	data, ok := img_response["data"].(map[string]interface{})
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Failed to get image data"))
		return ""
	}

	imagePath, ok := data["url"].(string)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Failed to get image url"))
		return ""
	}

	return imagePath
}

func ReadImage(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}