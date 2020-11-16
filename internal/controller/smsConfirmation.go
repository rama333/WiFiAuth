package controller

import (
	"WiFiAuth/internal/config"
	"WiFiAuth/internal/model"
	"WiFiAuth/internal/restapi"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func (c *Controller) SendCodeBySms(ctx *gin.Context) {

	country_code, err := strconv.Atoi(ctx.Request.FormValue("country-code"))
	if err != nil {
		restapi.ResponseBadRequest("Couldn't parse request body", ctx)
	}

	mobile_number, err := strconv.ParseInt(ctx.Request.FormValue("mobile_number"), 10, 64)
	if err != nil {
		restapi.ResponseBadRequest("Couldn't parse request body", ctx)
	}

	ip := ctx.ClientIP()
	log.Println(ip)

	mobileNumber := fmt.Sprint(country_code) + fmt.Sprint(mobile_number)

	smsCon := model.SMSConfirmationModel{Ip: ip, Code: model.GetCode(), CountAttempt: 0, TelNumber: mobileNumber}
	//s := model.NewSMSConfirmationModel(smsCon)

	if _, err := config.Config.REDISPOOL.Get().Do("SET", mobileNumber, smsCon); err != nil {
		log.Fatal(err)
	}

	//values, err := redis.String(config.Config.REDISPOOL.Get().Do("GET", mobileNumber))
	//if err != nil {
	//	fmt.Println("value rr")
	//	log.Fatal(err)
	//}

	stateCode := smsCon.SendCode()

	if stateCode == 201 {
		ctx.JSON(http.StatusOK, "ok")
	} else {
		restapi.ResponseInternalserverError("Failed to send message", ctx)
	}
}

func (с *Controller) CheckCodeByTelNumber(ctx *gin.Context) {

	mobile_number, err := strconv.ParseInt(ctx.Request.FormValue("mobile_number"), 10, 64)
	if err != nil {
		restapi.ResponseBadRequest("Couldn't parse request body", ctx)
	}

	mobileNumber := fmt.Sprint(mobile_number)

	code := ctx.Request.FormValue("code")

}
