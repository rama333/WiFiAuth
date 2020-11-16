package model

import (
	"WiFiAuth/internal/config"
	"bytes"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type SMSConfirmationModel struct {
	Ip           string `json:"ip"`
	Code         string `json:"code"`
	CountAttempt int    `json:"countAttempt"`
	TelNumber    string `json:"telNumber"`
}

type ReqForSendCode struct {
	Delivery_time string `json:"delivery_time"`
	Dest_addr     string `json:"dest_addr"`
	Id            int    `json:"id"`
	Sms_text      string `json:"sms_text"`
	Source_addr   string `json:"source_addr"`
}

func NewSMSConfirmationModel(model SMSConfirmationModel) *SMSConfirmationModel {
	return &model
}

func GetCode() string {
	charSet := "abcdedfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	rand.Seed(time.Now().Unix())

	var code strings.Builder

	for i := 0; i < 4; i++ {
		code.WriteString(string(charSet[rand.Intn(len(charSet)-1)]))
	}

	return code.String()
}

func (s *SMSConfirmationModel) SendCode() int {

	reqModel := ReqForSendCode{Dest_addr: s.TelNumber, Delivery_time: "", Id: 123, Sms_text: s.Code, Source_addr: "wi-fi"}

	url := "http://192.168.143.208:5000/sendMessage/api/v1/send"
	log.Println("URL:>", url)

	jsonStr, err := json.Marshal(reqModel)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.StatusCode)

	return resp.StatusCode
}

func CheckCode(mobilNumber string) {

	values, err := redis.String(config.Config.REDISPOOL.Get().Do("GET", mobilNumber))
	if err != nil {
		log.Println("value rr")
		log.Fatal(err)
	}

}
