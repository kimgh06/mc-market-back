package articles

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateArticle struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Head    string `json:"head"`
	CommentDisabled bool `json:"comment_disabled"`
	LikeDisabled    bool `json:"like_disabled"`
}

func updateArticle(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)
	
	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))	
		return
	}

	if user.ID != a.Queries.GetArticleAuthor(ctx, int64(id)) && user.Permissions != 2147483647 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to update this article"))
		return
	}

	body := UpdateArticle{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	if body.Head == "공지" {
		if user.Permissions != 2147483647 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to update article head"))
			return
		}
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

	fmt.Println(body.CommentDisabled)

	err = a.Queries.UpdateArticle(ctx, schema.UpdateArticleParams{
		ID:      int64(id),
		Title:   body.Title,
		Content: buffer.String(),
		Head:    sql.NullString{Valid: body.Head != "", String: body.Head},
		CommentDisabled: body.CommentDisabled,
		LikeDisabled: body.LikeDisabled,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}