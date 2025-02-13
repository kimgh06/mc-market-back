package routes

import (
	"maple/internal/api"
	"maple/internal/api/article_head"
	"maple/internal/api/article_likes"
	"maple/internal/api/articles"
	"maple/internal/api/banner"
	"maple/internal/api/comment"
	"maple/internal/api/payments"
	"maple/internal/api/products"
	"maple/internal/api/products_update_log"
	"maple/internal/api/products_versions"
	"maple/internal/api/status"
	"maple/internal/api/user"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/v1", func(a *api.MapleAPI) {
		user.InitializeRoutes(a)
		products.InitializeRoutes(a)
		payments.InitializeRoutes(a)
		articles.InitializeRoutes(a)
		comment.InitializeRoutes(a)
		article_likes.InitializeRoutes(a)
		article_head.InitializeRoutes(a)
		banner.InitializeRoutes(a)
		products_update_log.InitializeRoutes(a)
		products_versions.InitializeRoutes(a)
	})

	a.GET("/status/", status.Check)
}
