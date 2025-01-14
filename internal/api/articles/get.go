package articles

import (
	"database/sql"
	"errors"
	"maple/internal/api"
	"maple/internal/nullable"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)


type userGetResponse struct {
	ID       uint64  `json:"id"`
	Nickname *string `json:"nickname"`
}

type commentType struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id"`
	User 		userGetResponse `json:"user"`
	Replies  []commentType `json:"replies"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type getCommentsResponse struct {
	Comments []commentType `json:"comments"`
	Count    int64         `json:"count"`
}

func getReplies(ctx *gin.Context, commentID int64) []commentType {
	a:= api.Get(ctx)
	replies, err:= a.Queries.ListComments(ctx, schema.ListCommentsParams{ReplyFrom: &commentID})

	if err != nil {
		return []commentType{}
	}

	var response []commentType
	for _, reply := range replies {
		user, err := a.Queries.GetUserById(ctx, reply.UserID)
		if err != nil {
			return []commentType{}
		}

		response = append(response, commentType{
			ID:        strconv.FormatUint(uint64(reply.ID), 10),
			ArticleID: strconv.FormatUint(uint64(reply.ArticleID), 10),
			User: userGetResponse{
				ID:       uint64(reply.UserID),
				Nickname: func(ns *sql.NullString) *string {
					if ns.Valid {
						return &ns.String
					}
					return nil
				}(&user.Nickname),
			},
			Replies: getReplies(ctx, reply.ID),
			Content:   reply.Content,
			CreatedAt: reply.CreatedAt,
			UpdatedAt: reply.UpdatedAt,
		})
	}
	return response
}

func listComments(ctx *gin.Context, articleID uint64) getCommentsResponse {
	a := api.Get(ctx)

	comments, err := a.Queries.ListComments(ctx, schema.ListCommentsParams{ArticleID: articleID, PageID: 1})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return getCommentsResponse{}
	}
	
	if len(comments) == 0 {
		return getCommentsResponse{}
	}
	count, err := a.Queries.GetCountFromArticleId(ctx, schema.GetCountFromArticleId{ArticleID: articleID})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return getCommentsResponse{}
	}

	var response []commentType
	for _, comment := range comments {
		user, err := a.Queries.GetUserById(ctx, comment.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
			return getCommentsResponse{}
		}

		replies := getReplies(ctx, comment.ID)

		response = append(response, commentType{
			ID:        strconv.FormatUint(uint64(comment.ID), 10),
			ArticleID: strconv.FormatUint(uint64(comment.ArticleID), 10),
			User: userGetResponse{
				ID:       uint64(comment.UserID),
				Nickname: func(ns *sql.NullString) *string {
					if ns.Valid {
						return &ns.String
					}
					return nil
				}(&user.Nickname),
			},
			Replies: replies,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		})
	}

	return getCommentsResponse{
		Comments: response,
		Count:    count,
	}
}

type getArticleResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	Author    ArticleAuthor `json:"author"`
	Head      *string       `json:"head"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Views		  int64         `json:"views"`
	Comments	getCommentsResponse `json:"comments"`
	Likes int64 `json:"likes"`
	Dislikes int64 `json:"dislikes"`
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

	comments := listComments(ctx, uint64(article.Article.ID))

	ctx.JSON(http.StatusOK, getArticleResponse{
		ID:      strconv.FormatUint(uint64(article.Article.ID), 10),
		Title:   article.Article.Title,
		Content: article.Article.Content,
		Author: ArticleAuthor{
			ID:       strconv.FormatInt(article.User.ID, 10),
			Username: usernames[0],
			Nickname: article.User.Nickname.String,
		},
		Head:      nullable.StringToPointer(article.Article.Head),
		CreatedAt: article.Article.CreatedAt,
		UpdatedAt: article.Article.UpdatedAt,
		Views:    article.Article.Views,
		Comments: comments,
		Likes: article.Likes,
		Dislikes: article.DisLikes,
	})
}
