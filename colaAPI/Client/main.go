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
	token, err := GetToken(confYaml.APIServer)
	if err != nil || len(token) == 0 {
		token, err = ColaLogin(confYaml.APIServer)
		if err != nil || len(token) == 0 {
			// fmt.Println(err)
			fmt.Println("1,e")
			return
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
			fmt.Println("1,e")
			return
		}
		var s bool = false
		for {
			s = AddAccount(confYaml.APIServer, orderID, t)
			if s {
				break
			}
			time.Sleep(500)
		}
		if s {
			fmt.Println("0," + orderID)
			return
		}
	}
	if a != "0" {
		status := TowOrder(token, qrurl, a)
		if !status {
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


func ColaLogin(url string) (token string, err error) {
	URL := strings.Join([]string{url, "ColaLogin"}, "/")
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
