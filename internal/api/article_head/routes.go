package article_head

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/article_head/", func(a *api.MapleAPI) {
		a.GET("/list/", getHeadList)
		a.GET("/byid/:id/", getHeadById)
		a.GET("/byname/:name/", getHeadByName)
	})

	a.Route("/article_head/", func(a *api.MapleAPI) {
		a.Use(middlewares.RequireAuthentication(a))

		// add
		a.POST("/", createHead)
		// delete
		a.POST("/delete/:id/", deleteHead)
		// update
		a.POST("/update/:id/", updateHead)
	})
}
