package permissions

import (
	"github.com/gin-gonic/gin"
	"maple/internal/perrors"
	"maple/internal/schema"
	"net/http"
)

func CheckUserPermission(value UserPermission, flag UserPermission) bool {
	return value&flag != 0
}

func CheckUserPermissionCtx(ctx *gin.Context, user *schema.User, required UserPermission) bool {
	return CheckUserPermission(UserPermission(user.Permissions), required)
}

func RequireUserPermission(ctx *gin.Context, user *schema.User, required UserPermission) bool {
	if !CheckUserPermission(UserPermission(user.Permissions), required) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, perrors.InsufficientUserPermission.MakeJSON())
		return false
	}
	return true
}
