package model

import (
	"WiFiAuth/internal/config"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
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
	AnyNumber    string `json:"anyNumber"`
}

type ReqForSendCode struct {
	Delivery_time string `json:"delivery_time"`
	Dest_addr     string `json:"dest_addr"`
	Id            int    `json:"id"`
	Sms_text      string `json:"sms_text"`
	Source_addr   string `json:"source_addr"`
}

type SOAPResponse struct {
	XMLName xml.Name
	Body    struct {
		XMLName              xml.Name
		SessionQueryResponse struct {
			SessionQueryResponse struct {
				XMLName      xml.Name
				Result       int    `xml:"Result"`
				Description  string `xml:"Description"`
				SubscriberID string `xml:"SubscriberID"`
				Services     string `xml:"Services"`
				Quotas       string `xml:"Quotas"`
				IpAddress    string `xml:"IpAddress"`
				MacAddress   string `xml:"MacAddress"`
			} `xml:"SessionQueryResponse"`
		} `xml:"SessionQueryResponse"`
	}
}

type SOAPResponseLogon struct {
	XMLName xml.Name
	Body    struct {
		XMLName              xml.Name
		AccountLogonResponse struct {
			AccountLogonResponse struct {
				XMLName      xml.Name
				Result       int    `xml:"Result"`
				Description  string `xml:"Description"`
				SubscriberID string `xml:"SubscriberID"`
				Services     string `xml:"Services"`
				Quotas       string `xml:"Quotas"`
				IpAddress    string `xml:"IP_Address"`
				MacAddress   string `xml:"MAC_Address"`
			} `xml:"AccountLogonResponse"`
		} `xml:"AccountLogonResponse"`
	}
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"passwd"`
}

type Number struct {
	Id     string    `db:"id"`
	Number string    `db:"number"`
	State  string    `db:"state"`
	Date   time.Time `db:"date"`
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

func (s *SMSConfirmationModel) CallPhone(dst string, any string) int {

	url := "http://192.168.180.228:9052/api/sip/callOrder?dst=" + dst + "&ani=" + any
	log.Println("URL:>", url)

	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
		return 0
	}

	return resp.StatusCode
}

func (s *SMSConfirmationModel) GetNumber() string {

	number := []Number{}
	date := time.Now()
	date.Format("2006-01-02 15:04:05")

	for {
		err := config.Config.POSTGRESDB.Select(&number, "SELECT * from call_number where id = (select NEXTVAL('id_seq')) and state != 1")
		if err != nil || len(number) < 1 {
			fmt.Println(err)
			continue
		}
		break
	}

	_, err := config.Config.POSTGRESDB.Exec("UPDATE call_number SET state='1', date=$1 WHERE number=$2", date, number[0].Number)
	if err != nil {
		log.Println(err)
	}

	return number[0].Number
}

func UpdateStateCallNumber() {
	_, err := config.Config.POSTGRESDB.Exec("UPDATE call_number SET state='0' where date < CURRENT_TIMESTAMP - INTERVAL '360 second' and state != 0")
	if err != nil {
		log.Println(err)
	}
}

func ClearStateDB(number string) {
	_, err := config.Config.POSTGRESDB.Exec("UPDATE call_number SET state='0' WHERE number=$1", number)
	if err != nil {
		log.Println("err", err)
	}
}

func SetLogsSendSmsPhone(anynumber string, number string, ip string, port string, user_agent string) {
	date := time.Now()
	date.Format("2006-01-02 15:04:05")
	tx := config.Config.POSTGRESDB.MustBegin()
	tx.MustExec("INSERT INTO send_sms_phone_logs (date_added, anynumber, number, ip, port, user_agent) VALUES ($1, $2, $3,$4,$5,$6)", date, anynumber, number, ip, port, user_agent)
	tx.Commit()
}

func SetLogsBase(anynumber string, number string, ip string, port string, user_agent string, login string, password string, mac string) {
	date := time.Now()
	date.Format("2006-01-02 15:04:05")
	tx := config.Config.POSTGRESDB.MustBegin()
	tx.MustExec("INSERT INTO base_logs (date_added,login, password, mac, anynumber, number, ip, port, user_agent) VALUES ($1, $2, $3,$4,$5,$6,$7,$8,$9)", date, login, password, mac, anynumber, number, ip, port, user_agent)
	tx.Commit()
}

func SetLogsFailedIntroducesCode(anynumber string, number string, ip string, port string, user_agent string, failed_introduced_code string) {
	date := time.Now()
	date.Format("2006-01-02 15:04:05")
	tx := config.Config.POSTGRESDB.MustBegin()
	tx.MustExec("INSERT INTO failed_introduced_code_logs (date_added, failed_introduced_code, anynumber, number, ip, port, user_agent) VALUES ($1, $2, $3,$4,$5,$6,$7)", date, failed_introduced_code, anynumber, number, ip, port, user_agent)
	tx.Commit()
}

func GetUserAuth(numb string) (User, error) {
	url := "http://192.168.114.121:5001/billing/api/v1/set-msisdn"
	log.Println("URL:>", url)

	msisdnStruct := struct {
		Msisdn string `json:"msisdn"`
	}{}

	msisdnStruct.Msisdn = numb

	log.Println(msisdnStruct)

	jsonStr, err := json.Marshal(msisdnStruct)
	if err != nil {
		return User{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()

	var responseUser User

	err = json.NewDecoder(resp.Body).Decode(&responseUser)
	if err != nil {
		return User{}, err
	}

	log.Println("response Status:", resp.StatusCode)
	log.Println("response Status:", responseUser)

	return responseUser, nil
}

func CheckCode(mobilNumber string, code string) (SMSConfirmationModel, error) {

	var sms SMSConfirmationModel
	values, err := redis.String(config.Config.REDISPOOL.Get().Do("GET", mobilNumber))
	if err != nil {
		log.Println("value rr")
		return sms, err
	}

	log.Println(values)

	err = json.Unmarshal([]byte(values), &sms)
	if err != nil {
		return sms, err
	}

	log.Println(sms)

	return sms, nil
}

func AccountLogonRequest(ip string, port string, username string, password string) (int, string, string, error) {

	sessionQueryRequest := []byte(strings.TrimSpace(fmt.Sprintf(`
	<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
    <Body>
        <AccountLogonRequest xmlns="ws4p.irbis">
            <IP_Address xmlns="">%s</IP_Address>
            <IP_Port xmlns="">%s</IP_Port>
            <UserName xmlns="">%s</UserName>
            <Password xmlns="">%s</Password>
            <RememberMyDevice xmlns="">1</RememberMyDevice>
        </AccountLogonRequest>
    </Body>
</Envelope>`, ip, port, username, password)))

	req, err := http.NewRequest("POST", config.Config.SOAP_URL_WS4Portal, bytes.NewReader(sessionQueryRequest))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return -1, "", "", err
	}

	soapAction := "urn:AccountLogonRequest"

	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		//log.Fatal("Error on dispatching request. ", err.Error())
		return -1, "", "", err
	}

	log.Println(res.Status)

	result := new(SOAPResponseLogon)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		//log.Fatal("Error on unmarshaling xml. ", err.Error())
		return -1, "", "", err
	}

	log.Println(result)

	return result.Body.AccountLogonResponse.AccountLogonResponse.Result, result.Body.AccountLogonResponse.AccountLogonResponse.Description, result.Body.AccountLogonResponse.AccountLogonResponse.MacAddress, nil
}

func SessionQueryRequest(ip string, port string) (int, string, error) {
	//ip := "192.168.100.100"
	//port := "777"

	sessionQueryRequest := []byte(strings.TrimSpace(fmt.Sprintf(`
	<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
	<Body>
	<SessionQueryRequest xmlns="ws4p.irbis">
	<IP_Address xmlns="">%s</IP_Address>
	<IP_Port xmlns="">%s</IP_Port>
	</SessionQueryRequest>
	</Body>
	</Envelope>`, ip, port)))

	req, err := http.NewRequest("POST", config.Config.SOAP_URL_WS4Portal, bytes.NewReader(sessionQueryRequest))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return -1, "", err
	}

	soapAction := "urn:SessionQueryRequest"

	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		//log.Fatal("Error on dispatching request. ", err.Error())
		return -1, "", err
	}

	log.Println(res.Status)

	result := new(SOAPResponse)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		//log.Fatal("Error on unmarshaling xml. ", err.Error())
		return -1, "", err
	}

	log.Println(result)

	return result.Body.SessionQueryResponse.SessionQueryResponse.Result, result.Body.SessionQueryResponse.SessionQueryResponse.Description, nil

}
