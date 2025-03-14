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
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/godruoyi/go-snowflake"
)

type CreateArticle struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Head    string `json:"head"`
	CommentDisabled bool `json:"comment_disabled"`
	LikeDisabled    bool `json:"like_disabled"`
}

func createArticle(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	body := CreateArticle{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	if body.Head == "공지"{
		if user.Permissions != 2147483647 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to create article head"))
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

	
	ID := int64(snowflake.ID())

	row, err := a.Queries.CreateArticle(ctx, schema.CreateArticleParams{
		ID:      ID,
		Title:   body.Title,
		Content: buffer.String(),
		Author:  user.ID,
		Head:    sql.NullString{Valid: body.Head != "", String: body.Head},
		CommentDisabled: sql.NullBool{Valid: true, Bool: body.CommentDisabled},
		LikeDisabled: sql.NullBool{Valid: true, Bool: body.LikeDisabled},
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}
	
	headID, err := strconv.Atoi(row.Head.String)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.UnknownInternalError.MakeJSON(err.Error()))
		return
	}
	Head, err := a.Queries.GetArticleHeadByID(ctx, schema.GetArticleHeadByIDParams{ID: headID})
	if err != nil {
		fmt.Println(err, "No webhook URL")
		ctx.JSON(http.StatusOK, row)
		return
	}

	// Request HTTP POST to webhook_url
	if Head.WebhookURL != "" {
		http.Post(Head.WebhookURL, 
			"application/json", 
			bytes.NewBuffer([]byte(fmt.Sprintf(`{"content": "새로운 글이 등록되었습니다.\n https://mc-market.kr/articles/%d"}`, ID))))
	}
	
	ctx.JSON(http.StatusOK, row)
}
