package cmd

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"maple/internal/api"
	"maple/internal/conf"
	"maple/internal/middlewares"
	"maple/internal/routes"
)

var serveCommand = cobra.Command{
	Use:   "serve",
	Short: "Start maple and listen to requests",
	RunE:  handleServeCommand,
}

type logrusWriter struct{}

func (w logrusWriter) Write(data []byte) (n int, err error) {
	logrus.Infoln(string(data))
	return len(data), nil
}

func buildServeCommand() *cobra.Command {
	return &serveCommand
}

func handleServeCommand(cmd *cobra.Command, args []string) error {
	config, err := conf.LoadFromEnvironments()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration from environments\n")
		return nil
	} else {
		logrus.Println("Loaded configuration from environments")
	}

	if config.Logging.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debugln("Enabled debugging")
	}

	useMiddlewares := []gin.HandlerFunc{gin.Recovery(), middlewares.NewCORS()}
	if config.Logging.Requests {
		useMiddlewares = append(useMiddlewares, gin.Logger())
	}

	instance := api.NewMapleAPI(config, api.WithMiddlewares(useMiddlewares...))
	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			logrus.WithError(err).Errorln("Failed to close database connection in shutdown")
		}
	}(instance.Conn)

	if config.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	routes.InitializeRoutes(&instance)

	hosts := append([]string{config.API.Host}, config.API.AdditionalHosts...)

	logrus.Infof("Listening for requests on %v", hosts)
	return instance.Router.Run(hosts...)
}
