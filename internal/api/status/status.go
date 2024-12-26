package status

import (
	"github.com/gin-gonic/gin"
	"maple/internal/api"
	"time"
)

type Status struct {
	Status  string        `json:"status"`
	Version string        `json:"version"`
	Uptime  time.Duration `json:"uptime"`
}

func Check(ctx *gin.Context) {
	ctx.JSON(200, Status{
		Status:  "available",
		Version: api.Version,
		Uptime:  time.Now().Sub(api.Get(ctx).StartedAt),
	})
}
