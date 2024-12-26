package routes

import (
	"maple/internal/api"
	"maple/internal/api/articles"
	"maple/internal/api/payments"
	"maple/internal/api/products"
	"maple/internal/api/status"
	"maple/internal/api/user"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/v1", func(a *api.MapleAPI) {
		user.InitializeRoutes(a)
		products.InitializeRoutes(a)
		payments.InitializeRoutes(a)
		articles.InitializeRoutes(a)
	})

	a.GET("/status/", status.Check)
}
