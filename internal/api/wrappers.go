package api

import (
	"github.com/gin-gonic/gin"
)

func (a *MapleAPI) Handle(httpMethod string, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return a.Group.Handle(httpMethod, relativePath, handlers...)
}

func (a *MapleAPI) POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return a.Group.POST(relativePath, handlers...)
}

func (a *MapleAPI) GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return a.Group.GET(relativePath, handlers...)
}

func (a *MapleAPI) DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return a.Group.DELETE(relativePath, handlers...)
}

func (a *MapleAPI) PATCH(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return a.Group.PATCH(relativePath, handlers...)
}

func (a *MapleAPI) PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return a.Group.PUT(relativePath, handlers...)
}

func (a *MapleAPI) OPTIONS(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return a.Group.OPTIONS(relativePath, handlers...)
}

func (a *MapleAPI) HEAD(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return a.Group.HEAD(relativePath, handlers...)
}

func (a *MapleAPI) Any(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {

	return a.Group.Any(relativePath, handlers...)
}

func (a *MapleAPI) Route(relativePath string, fn func(api *MapleAPI), handlers ...gin.HandlerFunc) {
	copied := *a
	copied.Group = a.Group.Group(relativePath, handlers...)
	fn(&copied)
}

func (a *MapleAPI) Use(middlewares ...gin.HandlerFunc) {
	a.Group.Use(middlewares...)
}

func (a *MapleAPI) UseGlobal(middlewares ...gin.HandlerFunc) {
	a.Router.Use(middlewares...)
}

func (a *MapleAPI) With(middlewares ...gin.HandlerFunc) *gin.RouterGroup {
	group := a.Group.Group("/")
	group.Use(middlewares...)
	return group
}
