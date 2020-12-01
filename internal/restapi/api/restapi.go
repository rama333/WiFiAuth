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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func New(logger *zap.SugaredLogger) *RESTAPI {

	controller := controller.NewController()
	r := gin.New()
	r.Use(CORSMiddleware())
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("SendCodeBySms", controller.SendCodeBySms)
		v1.POST("CheckCode", controller.CheckCodeByTelNumber)
		v1.POST("CheckSession", controller.CheckSession)

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
