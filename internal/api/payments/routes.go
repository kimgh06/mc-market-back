package payments

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/payments", func(a *api.MapleAPI) {
		a.Use(middlewares.RequireAuthentication(a))

		a.GET("/", getPaymentList)
		a.POST("/", createPayment)
		a.POST("/:orderId/approve/", approvePayment)
		a.GET("/:orderId/", getPayment)
	})
}
