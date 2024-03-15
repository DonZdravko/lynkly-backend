package main

import (
	"fmt"
	"lynkly-backend/internal/logging"
	"lynkly-backend/internal/servers"
)

const (
	serviceName = "lynkly"
)

func main() {
	// Load configuration
	//mongoConfig := config.NewMongoConfig()
	//serverConfig := config.NewServerConfig()
	logger := logging.NewLogger(serviceName)
	logger.SetLevel(logging.DebugLevel)
	logger.Info(fmt.Sprintf("Starting %s...", serviceName))
	logger.Debug("Debug main.go")

	server := servers.NewUrlShortenerServer(fmt.Sprintf("127.0.0.1:%d", 18080), servers.ServerParams{
		Logger:     logger,
		ServiceUrl: "http://127.0.0.1:18080",
	})

	//// Start server
	err := server.Run()
	if err != nil {
		logger.Panic("Error encountered on running the server", "error", err)
	}

	//
	//// Start server
	//log.Printf("Server listening on port %s", serverConfig.Port)
	//log.Fatal(http.ListenAndServeTLS(":"+mongoConfig.ServerPort, mongoConfig.TLSCertFile, mongoConfig.TLSKeyFile, server.Router))
}
