package articles

import (
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"maple/internal/utilities"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type listArticlesElement struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Author    ArticleAuthor `json:"author"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Views 		int64         `json:"views"`
	HasImg 	bool          `json:"has_img"`
	CommentCount int64      `json:"comment_count"`
	Likes int64 `json:"likes"`
	Head 		*string       `json:"head"`
}

func listArticles(ctx *gin.Context) {
	a := api.Get(ctx)

	offsetQuery := ctx.DefaultQuery("offset", "0")
	sizeQuery := ctx.DefaultQuery("size", "20")

	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidQuery.MakeJSON(err.Error()))
		return
	}
	size, err := strconv.Atoi(sizeQuery)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidQuery.MakeJSON(err.Error()))
		return
	}
	
	size = utilities.Clamp(size, 0, 20)

	var rows []*schema.ListArticlesRow

	headId := ctx.Query("head_id")

	if headId == "0" {
		rows, err = a.Queries.ListArticles(ctx, schema.ListArticlesParams{
			Offset: int32(offset),
			Limit:  int32(size),
		})
	}else{ 
		headIDInt, err := strconv.Atoi(headId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidQuery.MakeJSON(err.Error()))
			return
		}
		rows, err = a.Queries.ListArticlesByHead(ctx, schema.ListArticlesByHeadParams{
			Head:   strconv.Itoa(headIDInt),
			Offset: int32(offset),
			Limit:  int32(size),
		})
	}
	

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.WithJSON(err.Error()))
		return
	}

	userIds := utilities.Map(rows, func(t *schema.ListArticlesRow) uint64 {
		return uint64(t.User.ID)
	})

	usernames, err := a.SurgeAPI.ResolveUsernamesAsMap(userIds)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return
	}

	articles := make([]listArticlesElement, len(rows))
	
	for i, row := range rows {
		articles[i] = listArticlesElement{
			ID:    strconv.FormatUint(uint64(row.Article.ID), 10),
			Title: row.Article.Title,
			Author: ArticleAuthor{
				ID:       strconv.FormatInt(row.User.ID, 10),
				Username: usernames[uint64(row.User.ID)],
				Nickname: row.User.Nickname.String,
			},
			CreatedAt: row.Article.CreatedAt,
			UpdatedAt: row.Article.UpdatedAt,
			Views:     row.Article.Views,
			HasImg:    row.HasImg,
			CommentCount: row.CommentCount,
			Likes: row.Likes,
			Head: &row.Article.Head.String,
		}
	}

	ctx.JSON(http.StatusOK, articles)
}
