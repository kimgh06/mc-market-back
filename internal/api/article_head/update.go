package article_head

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UpdateArticleHead struct {
	Name       string `json:"name"`
	IsAdmin    bool   `json:"is_admin"`
	WebhookURL string `json:"webhook_url"`
}

func updateHead(ctx *gin.Context) {
	a := api.Get(ctx)
	// check user permission
	user := middlewares.GetUser(ctx)
	if user.Permissions != 2147483647 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to create article head"))
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidParameter.MakeJSON("Invalid id parameter"))
		return
	}

	body := UpdateArticleHead{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	err = a.Queries.UpdateArticleHead(ctx, schema.ArticleHead{
		ID:         id,
		IsAdmin:    body.IsAdmin,
		Name:       body.Name,
		WebhookURL: body.WebhookURL,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}
