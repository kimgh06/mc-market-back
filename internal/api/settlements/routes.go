package settlements

import (
	"maple/internal/api"
	"maple/internal/middlewares"
)

func InitializeRoutes(a *api.MapleAPI) {
	a.Route("/settlements", func(a *api.MapleAPI) {
		a.Route("", func(a *api.MapleAPI) {
			a.Use(middlewares.RequireAuthentication(a))
			a.GET("/range", getSettlementsByUserIdAndMonths)
			a.GET("", getSettlementByUserIdAndMonth)
		})
	})
}
