package comment

import (
	"bytes"
	"database/sql"
	"html/template"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/godruoyi/go-snowflake"
)

type CreateComment struct {
	ReplyTo   int64 `json:"reply_to" binding:"omitempty"`
	Content   string `json:"content" binding:"required"`
}

func createComment(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	var body CreateComment
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	tmpl, err := template.New("comment.content").Parse(body.Content)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.UnknownInternalError.MakeJSON(err.Error()))
		return
	}

	articleIDStr := ctx.Param("article_id")
	if articleIDStr == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("article_id is required"))
		return
	}

	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("article_id must be an integer"))
		return
	}

	var buffer bytes.Buffer
	if err = tmpl.Execute(&buffer, body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidHTML.MakeJSON(err.Error()))
		return
	}

	commentID := int64(snowflake.ID())
	_, err = a.Queries.CreateComment(ctx, schema.CreateCommentParams{
		ID:        commentID,
		ArticleID: articleID,
		ReplyTo:   sql.NullInt64{Int64: body.ReplyTo, Valid: body.ReplyTo != 0},
		Content:   buffer.String(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}
