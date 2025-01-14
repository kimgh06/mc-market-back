package comment

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/comment/:article_id", func(a *api.MapleAPI) {
		a.GET("/", listComments)
		a.GET("/:comment_id/", getoneComment)
	})

	a.Route("/comment/:article_id", func(a *api.MapleAPI) {
		a.Use(middlewares.RequireAuthentication(a))

		a.POST("/", createComment)
		a.DELETE("/:comment_id/", deleteComment)
	})
}
