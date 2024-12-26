package articles

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"maple/pkg/permissions"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/articles", func(a *api.MapleAPI) {
		a.GET("/count/", countArticles)
		a.GET("/:id/", getArticle)
		a.GET("/", listArticles)
	})

	a.Route("/articles", func(a *api.MapleAPI) {
		a.Use(middlewares.RequireAuthentication(a), middlewares.RequireUserPermission(permissions.ManageArticles))

		a.POST("/", createArticle)
		a.DELETE("/:id/", deleteArticle)
	})
}
