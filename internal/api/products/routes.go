package products

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/products", func(a *api.MapleAPI) {
		a.Route("/", func(a *api.MapleAPI) {
			a.Use(middlewares.RequireAuthentication(a))

			a.POST("/", createProduct)
			a.GET("/", listProducts)
		})

		a.Route("/", func(a *api.MapleAPI) {
			a.Use(middlewares.UseAuthentication(a))

			a.GET("/:id/", getProduct)
		})
	})
}
