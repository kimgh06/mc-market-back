package article_likes

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/article_likes/:article_id", func(a *api.MapleAPI) {
		a.GET("/", getArticleLikesAndDisLikes)
	})

	a.Route("/article_likes/:article_id", func(a *api.MapleAPI) {
		a.Use(middlewares.RequireAuthentication(a))

		a.POST("/", createArticleLike)
		a.DELETE("/", deleteArticleLike)
	})
}
