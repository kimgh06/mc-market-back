package api

import (
	"crypto/tls"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"maple/internal/conf"
	"maple/internal/schema"
	"maple/internal/storage"
	"net/http"
)

type MapleAPI struct {
	Router *gin.Engine
	Group  *gin.RouterGroup

	Conn    *sql.DB
	Queries *schema.Queries
	Config  *conf.MapleConfigurations

	SurgeHTTP *http.Client

	JWKS *MapleJWKS
}

type MapleAPIOptions func(*MapleAPI)

func WithMiddlewares(middlewares ...gin.HandlerFunc) MapleAPIOptions {
	return func(api *MapleAPI) {
		api.Router.Use(middlewares...)
		api.Group = api.Router.Group("/")
	}
}

func NewMapleAPI(config *conf.MapleConfigurations, options ...MapleAPIOptions) MapleAPI {
	connection := storage.CreateDatabaseConnection(&config.Database)

	router := gin.New()
	router.RedirectTrailingSlash = false

	jwks, err := NewMapleJWKS(config)
	if err != nil {
		logrus.Errorln(err)
		logrus.
			WithField("component", "api").
			Warnln("There was an error creating JWKs manager. This will disable most of features")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	instance := MapleAPI{
		Router: router,
		Group:  router.Group("/"),

		Conn:    connection,
		Queries: schema.New(connection),
		Config:  config,

		SurgeHTTP: &http.Client{Transport: tr},

		JWKS: jwks,
	}

	router.Use(func(context *gin.Context) {
		context.Set(ContextAPIKey.String(), &instance)
	})

	for o := range options {
		opt := options[o]
		opt(&instance)
	}

	return instance
}
