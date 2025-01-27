package article_head

import (
	"fmt"
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/internal/perrors"
	"maple/internal/schema"
	"math/rand/v2"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateArticleHead struct {
	Name  string `json:"name"`
	IsAdmin bool `json:"is_admin"`
}

func createHead(ctx *gin.Context) {
	a := api.Get(ctx)
	// check user permission
	user := middlewares.GetUser(ctx)
	if user.Permissions != 2147483647 {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON("You don't have permission to create article head"))
		return
	}
	
	body := CreateArticleHead{}
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, perrors.InvalidJSON.MakeJSON(err.Error()))
		return
	}

	headID := rand.IntN(100000)
	fmt.Println("headID: ", headID, "name: ", body.Name)
	
	err := a.Queries.CreateArticleHead(ctx, schema.ArticleHead{
		ID:   headID,
		IsAdmin: body.IsAdmin,
		Name:  body.Name,
	})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, perrors.FailedDatabase.MakeJSON(err.Error()))
		return
	}

	ctx.Status(http.StatusCreated)
}
