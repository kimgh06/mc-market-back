package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maple/internal/perrors"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
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
		fmt.Println(body, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Invalid response format: "+err.Error()))
		return ""
	}
	
	fmt.Println(img_response)
	imagePath, ok := img_response["filename"].(string)
	if !ok {
		fmt.Println(imagePath, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON("Failed to get image data"))
		return ""
	}

	return url+"/images/"+imagePath
}

func ReadImage(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	// Resize to 30% of original size
	bounds := img.Bounds()
	resize := 0.5
	newWidth := int(float64(bounds.Max.X) * resize)
	newHeight := int(float64(bounds.Max.Y) * resize)
	resizedImg := imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)

	// Convert back to bytes
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, resizedImg, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}