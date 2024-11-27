package routes

import (
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"maple/internal/api/products"
	"maple/internal/api/user"
	"net/http"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/v1", func(a *api.MapleAPI) {
		user.InitializeRoutes(a)
		products.InitializeRoutes(a)
	})

	a.GET("/", func(context *gin.Context) {
		context.JSON(http.StatusOK, "안녕.")
	})
}
