package products_update_log

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/products_log/", func(a *api.MapleAPI) {

		a.GET("/list/:product_id/", getUpdateList)
		a.GET("/:id/", getUpdateOne)
	})
	
	a.Route("/products_log/", func(a *api.MapleAPI) {
		a.Use(middlewares.RequireAuthentication(a))
		
		a.POST("/create/:product_id/", createProductUpdateLog)
		a.POST("/update/:id/", update)
		a.POST("/delete/:id/", deleteLog)
	})
}
