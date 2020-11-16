package main

import (
	"WiFiAuth/internal/config"
	"WiFiAuth/internal/resources"
	"WiFiAuth/internal/restapi/api"
	"fmt"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	slogger := logger.Sugar()
	slogger.Info("Starting the application...")
	slogger.Info("Reading configuration and initializing resources...")

	if err := config.LoadConfig("/smsService"); err != nil {
		panic(fmt.Errorf("invalid application configuration: %s", err))
	}

	rsc, err := resources.New(slogger)
	if err != nil {
		slogger.Fatalw("Can't initialize resources.", "err", rsc)
	}

	slogger.Info("Configuring the application units...")

	//diag := diagnostics.New(slogger, config.Config.DIAGPORT, rsc.Healthz)
	//diag.Start(slogger)
	//slogger.Info("The application is ready to serve requests.")

	rapi := api.New(slogger)
	rapi.Start(config.Config.RESTAPIPort)

}
