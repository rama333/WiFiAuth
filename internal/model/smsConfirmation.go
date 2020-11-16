package model

import (
	"math/rand"
	"strings"
	"time"
)

type SMSConfirmationModel struct {
	Ip string `json:"ip"`
	Code string `json:"code"`
	CountAttempt int `json:"countAttempt"`
}




func GetCode() (string) {
	charSet := "abcdedfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	rand.Seed(time.Now().Unix())

	var code strings.Builder

	for i := 0; i < 4; i++ {
		code.WriteString(string(charSet[rand.Intn(len(charSet)-1)]))
	}

	return code.String()
}
