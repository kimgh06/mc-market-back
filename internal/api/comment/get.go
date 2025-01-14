package comment

import (
	"database/sql"
	"errors"
	"fmt"
	"maple/internal/api"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type userGetResponse struct {
	ID       uint64  `json:"id"`
	Nickname *string `json:"nickname"`
}

type getCommentResponse struct {
	ID        string    `json:"id"`
	ArticleID string    `json:"article_id"`
	User 		userGetResponse `json:"user"`
	Replies  []getCommentResponse `json:"replies"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func getReplies(ctx *gin.Context, commentID int64) []getCommentResponse {
	a:= api.Get(ctx)
	replies, err:= a.Queries.ListComments(ctx, schema.ListCommentsParams{ReplyFrom: &commentID})

	if err != nil {
		return []getCommentResponse{}
	}

	var response []getCommentResponse
	for _, reply := range replies {
		user, err := a.Queries.GetUserById(ctx, reply.UserID)
		if err != nil {
			return []getCommentResponse{}
		}

		response = append(response, getCommentResponse{
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

func listComments(ctx *gin.Context){
	a := api.Get(ctx)

	articleID, err := api.GetUint64FromParam(ctx, "article_id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake)
		return
	}

	page  := strings.Replace(ctx.Query("page"), "/", "", 1)
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		// Handle error if conversion fails
		fmt.Println("Error converting to int:", err)
		return
	}

	comments, err := a.Queries.ListComments(ctx, schema.ListCommentsParams{ArticleID: articleID, PageID: int32(pageInt)})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}
	
	if len(comments) == 0 {
		ctx.JSON(http.StatusOK, []getCommentResponse{})
		return
	}
	count, err := a.Queries.GetCountFromArticleId(ctx, schema.GetCountFromArticleId{ArticleID: articleID})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	var response []getCommentResponse
	for _, comment := range comments {
		user, err := a.Queries.GetUserById(ctx, comment.UserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
			return
		}

		replies := getReplies(ctx, comment.ID)

		response = append(response, getCommentResponse{
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

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"comments": response,
		"count":    count,
	})
}

func getoneComment(ctx *gin.Context) {
	a := api.Get(ctx)

	commentID, err := api.GetUint64FromParam(ctx, "comment_id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake)
		return
	}

	comment, err := a.Queries.GetComment(ctx, int64(commentID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.CommentNotFound.MakeJSON())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	user, err := a.Queries.GetUserById(ctx, comment.UserID)

	ctx.JSON(http.StatusOK, getCommentResponse{
		ID:        strconv.FormatUint(uint64(comment.ID), 10),
		ArticleID: strconv.FormatUint(uint64(comment.ArticleID), 10),
		User: userGetResponse{
			ID:       uint64(comment.UserID),
			Nickname: func() *string {
				if user.Nickname.Valid {
					return &user.Nickname.String
				}
				return nil
			}(),
		},
		Replies: getReplies(ctx, comment.ID),
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	})
}
