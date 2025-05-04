package adcard

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/adcard", func(a *api.MapleAPI) {
		a.Route("/", func(a *api.MapleAPI) {
			a.Use(middlewares.RequireAuthentication(a))

			a.POST("/upload/", uploadImage)
			a.POST("/delete/:id/", deleteImage)
			a.POST("/update/:id/", updateImage)
		})
		a.GET("/list/", getListImage)
		a.GET("/image/:path", getImageFromUrl)
		a.GET("/image/:path/", getImageFromUrl)
		a.GET("/:id/", getImage)
	})
}
