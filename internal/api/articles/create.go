package articles

import (
	"bytes"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/godruoyi/go-snowflake"
	"html/template"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
)

type CreateArticle struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Head    string `json:"head"`
}

func createArticle(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	body := CreateArticle{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	tmpl, err := template.New("article.content").Parse(body.Content)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.UnknownInternalError.MakeJSON(err.Error()))
		return
	}
	buffer := new(bytes.Buffer)

	if err = tmpl.Execute(buffer, body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidHTML.MakeJSON(err.Error()))
		return
	}

	_, err = a.Queries.CreateArticle(ctx, schema.CreateArticleParams{
		ID:      int64(snowflake.ID()),
		Title:   body.Title,
		Content: buffer.String(),
		Author:  user.ID,
		Head:    sql.NullString{Valid: body.Head != "", String: body.Head},
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}
