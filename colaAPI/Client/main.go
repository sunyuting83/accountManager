package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"gopkg.in/yaml.v2"
)

type Config struct {
	APIServer string `yaml:"APIServer"`
}

type GetApiRequest struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type CreateRequest struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   Data   `json:"data"`
}
type ColaAccountRequest struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type TwoAuthRequest struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   Data   `json:"data"`
}

type Wxdata struct {
	Wxid     interface{} `json:"wxid"`
	Nickname string      `json:"nickname"`
}
type Wxid struct {
	Wxid     string `json:"wxid"`
	Nickname string `json:"nickname"`
}
type Order struct {
	ID            string `json:"id"`
	Oid           string `json:"oid"`
	Ordermode     int    `json:"ordermode"`
	Ordermodetext string `json:"ordermodetext"`
}
type Data struct {
	Fullurl string `json:"fullurl"`
	Wxdata  Wxdata `json:"wxdata"`
	Wxid    Wxid   `json:"wxid"`
	Order   Order  `json:"order"`
}

type RequestLogin struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		MemberID      string `json:"member_id"`
		MemberUser    string `json:"member_user"`
		AgentID       string `json:"agent_id"`
		MemberBalance string `json:"member_balance"`
		MemberState   string `json:"member_state"`
	} `json:"data"`
	Token string `json:"token"`
}

func main() {
	var (
		f string
		s string
		a string
		p string
		t string
	)
	flag.StringVar(&f, "f", "0", "file")
	flag.StringVar(&s, "s", "0", "first")
	flag.StringVar(&a, "a", "0", "account")
	flag.StringVar(&p, "p", "2370", "projectid")
	flag.StringVar(&t, "t", "3", "status")
	flag.Parse()

	OS := runtime.GOOS

	CurrentPath, _ := GetCurrentPath()
	confYaml, err := CheckConfig(OS, CurrentPath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	// get token from API Server
	token, err := GetToken(confYaml.APIServer)
	if err != nil || len(token) == 0 {
		// haven't token, so create Login for cola
		var colaRequest ColaAccountRequest
		LoginStatus := false
		for {
			// get cola username and password an token
			// if has a active at server, wait it done.
			colaRequest, _ = GetColaAccount(confYaml.APIServer)
			if colaRequest.Status == 0 {
				LoginStatus = true
				token = colaRequest.Token
				break
			}
			time.Sleep(time.Duration(10) * time.Second)
		}
		if LoginStatus && len(token) == 0 {
			token, err = ColaLogin(colaRequest)
			if err != nil || len(token) == 0 {
				SetColaToken(confYaml.APIServer, token)
				fmt.Println("1,e")
				return
			}
			SetColaToken(confYaml.APIServer, token)
		}
	}
	// fmt.Println(token)
	qrurl, err := GetPaymentStr(f)
	if err != nil {
		fmt.Println("1,e")
		return
	}
	if s == "1" {
		status, orderID := CreateOrder(token, qrurl, p)
		if !status {
			SetColaToken(confYaml.APIServer, "")
			fmt.Println("1,e")
			return
		}
		var sa bool = false
		for {
			sa = AddAccount(confYaml.APIServer, orderID, t)
			if sa {
				break
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
		if sa {
			fmt.Println("0," + orderID)
			return
		}
	}
	if a != "0" {
		status := TowOrder(token, qrurl, a)
		if !status {
			SetColaToken(confYaml.APIServer, "")
			fmt.Println("1,e")
			return
		}
		fmt.Println("0,o")
		return
	}
}

func TowOrder(token, uri, a string) (status bool) {
	Params := strings.Join([]string{"and_id=2&o_id=", a, "&url=", uri}, "")
	URL := "http://tiancaiapi.tablecando.cn/api/Order/twoauth"
	req, _ := http.NewRequest("POST", URL, strings.NewReader(Params))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", token)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)
	var request *TwoAuthRequest
	json.Unmarshal(respByte, &request)
	status = false
	if request.Status == 205 {
		status = true
	}
	return
}

func CreateOrder(token, uri, p string) (status bool, orderid string) {
	Params := strings.Join([]string{"type=99&projectid=", p, "&url=", uri}, "")
	URL := "http://tiancaiapi.tablecando.cn/api/Order/CreateOrder"
	req, _ := http.NewRequest("POST", URL, strings.NewReader(Params))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", token)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	if err != nil {
		// fmt.Println(err)
		return false, ""
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)
	var request *CreateRequest
	json.Unmarshal(respByte, &request)
	status = false
	orderid = ""
	// fmt.Println(request)
	if request.Status == 200 {
		status = true
		orderid = request.Data.Order.ID
	}
	return
}

func GetToken(url string) (token string, err error) {
	URL := strings.Join([]string{url, "GetColaToken"}, "/")
	// fmt.Println(URL)
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)
	var request GetApiRequest
	json.Unmarshal(respByte, &request)
	if request.Status == 1 {
		return "", errors.New(request.Message)
	}
	return request.Token, nil
}

func GetColaAccount(url string) (request ColaAccountRequest, err error) {
	URL := strings.Join([]string{url, "ColaAccount"}, "/")
	// fmt.Println(URL)
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	if err != nil {
		return request, err
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respByte, &request)
	if request.Status == 1 {
		return request, errors.New(request.Message)
	}
	return request, nil
}
func SetColaToken(url, token string) {
	URL := strings.Join([]string{url, "SetColaToken"}, "/")
	Params := strings.Join([]string{"token=", token}, "")
	req, _ := http.NewRequest("PUT", URL, strings.NewReader(Params))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")
	_, _ = (&http.Client{Timeout: 35 * time.Second}).Do(req)
}

func ColaLogin(colaRequest ColaAccountRequest) (token string, err error) {
	u := "Member/login"
	url := strings.Join([]string{"http://tiancaiapi.tablecando.cn/api/", u}, "")

	Params := strings.Join([]string{"member_user=", colaRequest.UserName, "&member_pwd=", colaRequest.Password}, "")
	req, _ := http.NewRequest("POST", url, strings.NewReader(Params))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	var request *RequestLogin
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)

	json.Unmarshal(respByte, &request)
	if request.Status == "200" {
		token = request.Token
	}
	return
}

func AddAccount(url, oid, t string) bool {
	URL := strings.Join([]string{url, "AddAccount?account="}, "/")
	URL = strings.Join([]string{URL, oid, "&status=", t}, "")
	// fmt.Println(URL)
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)
	var request GetApiRequest
	json.Unmarshal(respByte, &request)
	status := true
	if request.Status == 1 {
		status = false
	}
	return status
}

// CheckConfig check config
func CheckConfig(OS, CurrentPath string) (conf *Config, err error) {
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	var confYaml *Config
	yamlFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	return confYaml, nil
}

// GetCurrentPath Get Current Path
func GetCurrentPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	return dir, nil
}

func GetPaymentStr(fi string) (paymentCodeUrl string, err error) {
	file, err := os.Open(fi)
	if err != nil {
		fmt.Println("a" + err.Error())
		return "", err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("b" + err.Error())
		return "", err
	}
	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		fmt.Println("c" + err.Error())
		return "", err
	}
	// fmt.Println(result.String())
	return result.String(), nil
}
