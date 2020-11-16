package api

import (
	"WiFiAuth/internal/controller"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RESTAPI struct {
	server *gin.Engine
	logger *zap.SugaredLogger
}

// @title sms service report Swagger API
// @version 1.0
// @description Swagger API for Golang Project sms service report Swagger API.
// @termsOfService http://swagger.io/terms/
// @host 192.168.114.145:8080
// @BasePath /api/v1

func New(logger *zap.SugaredLogger) *RESTAPI {

	controller := controller.NewController()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("StatMessage", controller.SendCodeBySms)

	}

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return &RESTAPI{server: r, logger: logger}
}

func (rapi *RESTAPI) Start(port int) {

	rapi.server.Run(fmt.Sprintf(":%v", port))

	//go func() {
	//	rapi.server.Run(fmt.Sprintf(":%v", port))
	//}()
}
