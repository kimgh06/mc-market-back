package comment

import (
	"database/sql"
	"errors"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/pkg/permissions"
	"net/http"

	"github.com/gin-gonic/gin"
)

func deleteComment(ctx *gin.Context) {
	a := api.Get(ctx)
	user := middlewares.GetUser(ctx)

	id, err := api.GetUint64FromParam(ctx, "comment_id")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidSnowflake.MakeJSON(err.Error()))
		return
	}

	comment, err := a.Queries.GetComment(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, perrors.CommentNotFound.MakeJSON(err.Error()))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if comment.UserID != user.ID && !permissions.CheckUserPermission(permissions.UserPermission(user.Permissions), permissions.ManageComments) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON())
		return
	}

	if err = a.Queries.DeleteReplies(ctx, int64(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	if err = a.Queries.DeleteComment(ctx, int64(id)); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusOK)
}
