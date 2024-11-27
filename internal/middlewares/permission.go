package middlewares

import (
	"github.com/gin-gonic/gin"
	"maple/internal/perrors"
	"maple/pkg/permissions"
	"net/http"
)

func RequireUserPermission(flag permissions.UserPermission) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := GetUser(c)

		if !permissions.CheckUserPermission(permissions.UserPermission(user.Permissions), flag) {
			c.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON())
			return
		}

		c.Next()
	}
}
