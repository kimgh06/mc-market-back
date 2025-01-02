package articles

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/perrors"
	"net/http"
	"strconv"
	"time"
)

type getArticleResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Author    ArticleAuthor `json:"author"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func getArticle(ctx *gin.Context) {
	a := api.Get(ctx)

	id, err := api.GetUint64FromParam(ctx, "id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake)
		return
	}

	article, err := a.Queries.GetArticle(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.ArticleNotFound.MakeJSON())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	usernames, err := a.SurgeAPI.ResolveUsernames([]uint64{uint64(article.User.ID)})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedAPI.MakeJSON(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, getArticleResponse{
		ID:      strconv.FormatUint(uint64(article.Article.ID), 10),
		Title:   article.Article.Title,
		Content: article.Article.Content,
		Author: ArticleAuthor{
			ID:       strconv.FormatInt(article.User.ID, 10),
			Username: usernames[0],
			Nickname: article.User.Nickname.String,
		},
		CreatedAt: article.Article.CreatedAt,
		UpdatedAt: article.Article.UpdatedAt,
	})
}
