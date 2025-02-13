package products_versions

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/products_versions/", func(a *api.MapleAPI) {

		a.GET("/list/:product_id/", getVersionList)
		a.GET("/:id/", getOneVersion)
	})
	
	a.Route("/products_versions/", func(a *api.MapleAPI) {
		a.Use(middlewares.RequireAuthentication(a))
		
		a.POST("/create/:product_id/", createVersion)
		a.POST("/update/:id/", update)
		a.POST("/delete/:id/", deleteVersion)
	})
}
