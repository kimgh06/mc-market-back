package api

import "github.com/gin-gonic/gin"

type contextKey string

func (c contextKey) String() string {
	return "maple.api.context-key " + string(c)
}

const (
	ContextAPIKey   = contextKey("api_instance")
	ContextTokenKey = contextKey("token")
	ContextClaimKey = contextKey("claim")
	ContextUserKey  = contextKey("user")
)

func Get(ctx *gin.Context) *MapleAPI {
	obj, _ := ctx.Get(ContextAPIKey.String())
	return obj.(*MapleAPI)
}
