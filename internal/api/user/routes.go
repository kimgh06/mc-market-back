package user

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/user", func(a *api.MapleAPI) {
		a.Route("/", func(a *api.MapleAPI) {
			a.Use(middlewares.RequireAuthentication(a))

			a.GET("/session/", getSessionUser)
			a.POST("/avatar/", uploadImage)
			a.GET("/revenues/", getUnclaimedRevenues)
			a.GET("/", listUsers)
			a.POST("/:id/", updateUser)
		})

		a.Route("/", func(a *api.MapleAPI) {
			a.Use(middlewares.UseAuthentication(a))

			a.GET("/:id/", getUser)
			a.GET("/:id/avatar/", getImage)
		})

		a.POST("/", createUser)
	})
}
