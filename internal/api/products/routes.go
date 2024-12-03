package products

import (
	"maple/internal/api"
	"maple/internal/middlewares"
	"path/filepath"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/products", func(a *api.MapleAPI) {
		a.Route("/", func(a *api.MapleAPI) {
			a.Use(middlewares.RequireAuthentication(a))

			a.POST("/", createProduct)
			a.POST("/:id/", updateProduct)
			a.DELETE("/:id/", deleteProduct)
			a.POST("/:id/image/", uploadImage)
		})

		a.Route("/", func(a *api.MapleAPI) {
			a.Use(middlewares.UseAuthentication(a))

			a.GET("/:id/", getProduct)
			a.GET("/", listProducts)
			a.GET("/:id/image/", getImage)
		})

		a.Route("/images/", func(a *api.MapleAPI) {
			a.Group.Static(filepath.Join(a.Config.Storage.ImagesPath, "products"), "/")
		})
	})
}
