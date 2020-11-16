package controller

import (
	"WiFiAuth/internal/config"
	"WiFiAuth/internal/model"
	"WiFiAuth/internal/restapi"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"strconv"
)

func(c * Controller) SendCodeBySms(ctx *gin.Context)  {

	country_code, err := strconv.Atoi(ctx.Request.FormValue("country-code"))
	if err != nil{
		restapi.ResponseBadRequest("Couldn't parse request body", c)
	}

	mobile_number, err := strconv.ParseInt(ctx.Request.FormValue("mobile_number"), 10, 64)
	if err != nil{
		restapi.ResponseBadRequest("Couldn't parse request body", c)
	}


	ip := ctx.ClientIP()

	country_code_string := strconv.Itoa(country_code)
	mobile_number_string := fmt.Sprint(mobile_number)

	smsCon := model.SMSConfirmationModel{Ip: ip, Code: model.GetCode(), CountAttempt: 0}

	if _, err := config.Config.REDISPOOL.Get().Do("SET", (country_code_string+mobile_number_string), smsCon); err != nil {
		log.Fatal(err)
	}


	values, err := redis.String(config.Config.REDISPOOL.Get().Do("GET", (country_code_string+mobile_number_string)))
	if err != nil {
		fmt.Println("value rr")
		log.Fatal(err)
	}


	ctx.JSON(http.StatusOK, values)


}