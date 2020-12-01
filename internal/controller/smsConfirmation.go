package controller

import (
	"WiFiAuth/internal/config"
	"WiFiAuth/internal/model"
	"WiFiAuth/internal/restapi"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"strconv"
)

func (c *Controller) SendCodeBySms(ctx *gin.Context) {

	confirModel := model.SMSConfirmationModel{}

	country_code, err := strconv.Atoi(ctx.Request.FormValue("country-code"))
	if err != nil {
		restapi.ResponseBadRequest("Couldn't parse request body", ctx)
	}

	mobile_number, err := strconv.ParseInt(ctx.Request.FormValue("mobile_number"), 10, 64)
	if err != nil {
		restapi.ResponseBadRequest("Couldn't parse request body", ctx)
	}

	ip, port, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		log.Println(err)
	}

	log.Println(port)

	mobileNumber := fmt.Sprint(country_code) + fmt.Sprint(mobile_number)

	anyNumber := confirModel.GetNumber()

	smsCon := model.SMSConfirmationModel{Ip: ip, Code: anyNumber[6:len(anyNumber)], CountAttempt: 0, TelNumber: mobileNumber, AnyNumber: anyNumber}

	model.SetLogsSendSmsPhone(anyNumber, mobileNumber, ip, port, ctx.Request.UserAgent())

	//model.NewSMSConfirmationModel(smsCo)

	sms, err := json.Marshal(smsCon)
	if err != nil {
		log.Println(err)
	}

	if _, err := config.Config.REDISPOOL.Get().Do("SET", mobileNumber, sms); err != nil {
		log.Fatal(err)
	}

	//values, err := redis.String(config.Config.REDISPOOL.Get().Do("GET", mobileNumber))
	//if err != nil {
	//	fmt.Println("value rr")
	//	log.Fatal(err)
	//}

	stateCode := smsCon.SendCode()
	stateMobileCode := smsCon.CallPhone(ctx.Request.FormValue("mobile_number"), anyNumber)

	log.Println("--------------")
	log.Println(stateMobileCode)
	log.Println("________________")

	if stateCode == 201 || stateMobileCode == 200 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":        "200",
			"error":       "0",
			"SendRequest": "1",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":        "200",
			"error":       "1",
			"SendRequest": "0",
		})
	}
}

func (с *Controller) CheckCodeByTelNumber(ctx *gin.Context) {

	ip, port, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		log.Println(err)
	}

	mobile_number, err := strconv.ParseInt(ctx.Request.FormValue("mobile_number"), 10, 64)
	if err != nil {
		restapi.ResponseBadRequest("Couldn't parse request body", ctx)
	}

	mobileNumber := fmt.Sprint(mobile_number)

	log.Println(mobileNumber)

	code := ctx.Request.FormValue("code")

	//var smsConfir model.SMSConfirmationModel

	smsConfir, err := model.CheckCode(mobileNumber, code)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":         "200",
			"error":        "1",
			"LogonRequest": "0",
		})
		return
	}

	if smsConfir.Code == code {
		var user model.User
		user, err = model.GetUserAuth(mobileNumber)

		log.Println("login", user.Login)
		log.Println("password", user.Password)
		log.Println("ip", ip)
		log.Println("port", port)

		log.Println("client ip", ctx.ClientIP())

		code, des, mac, err := model.AccountLogonRequest(ctx.ClientIP(), port, user.Login, user.Password)
		if err != nil {
			log.Println(err.Error() + des)
			ctx.JSON(http.StatusOK, gin.H{
				"code":         "200",
				"error":        "1",
				"LogonRequest": "0",
			})
		}

		if code == 0 {
			model.SetLogsBase(smsConfir.AnyNumber, smsConfir.TelNumber, ctx.ClientIP(), port, ctx.Request.UserAgent(), user.Login, user.Password, mac)
			ctx.JSON(http.StatusOK, gin.H{
				"code":         "200",
				"error":        "0",
				"LogonRequest": "1",
			})

			model.ClearStateDB(smsConfir.AnyNumber)

		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code":         "200",
				"error":        "1",
				"LogonRequest": "0",
			})
		}
	} else {

		ctx.JSON(http.StatusOK, gin.H{
			"code":         "200",
			"error":        "0",
			"LogonRequest": "0",
		})

		model.SetLogsFailedIntroducesCode(smsConfir.AnyNumber, mobileNumber, ip, port, ctx.Request.UserAgent(), code)
	}

}
func (c *Controller) CheckSession(ctx *gin.Context) {

	ip, port, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		log.Println(err)
	}

	log.Println(ip)
	log.Println(port)

	code, des, err := model.SessionQueryRequest(ip, port)

	if err != nil {
		log.Println(err.Error() + des)
		ctx.JSON(http.StatusOK, gin.H{
			"error": "1",
		})
	}

	if code == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"error":         "0",
			"activeSession": "1",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"error":         "0",
			"activeSession": "0",
		})
	}

}
