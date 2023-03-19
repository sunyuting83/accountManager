package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
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
	APIServer string `yaml:"Host"`
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
	)
	flag.StringVar(&f, "f", "0", "file")
	flag.StringVar(&s, "s", "0", "status")
	flag.StringVar(&a, "a", "0", "account")
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
	if err != nil {
		fmt.Println("1")
		return
	}
	qrurl, err := GetPaymentStr(f)
	if err != nil {
		fmt.Println("1")
		return
	}
	if s == "1" {
		status, orderID := CreateOrder(token, qrurl)
		if !status {
			fmt.Println("1")
			return
		}
		var s bool = false
		for {
			s = AddAccount(confYaml.APIServer, orderID)
			if s {
				break
			}
			time.Sleep(500)
		}
		if s {
			fmt.Println("0")
			return
		}
	}
	if a != "0" {
		status := TowOrder(token, qrurl, a)
		if !status {
			fmt.Println("1")
			return
		}
		fmt.Println("0")
		return
	}
	fmt.Println("1")
}

func TowOrder(token, uri, a string) (status bool) {
	Params := strings.Join([]string{"and_id=2&o_id", a, "&url=", uri}, "")
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
	if request.Status == 200 {
		status = true
	}
	return
}

func CreateOrder(token, uri string) (status bool, orderid string) {
	Params := strings.Join([]string{"type=99&projectid=2068&url", uri}, "")
	URL := "http://tiancaiapi.tablecando.cn/api/Order/CreateOrder"
	req, _ := http.NewRequest("POST", URL, strings.NewReader(Params))
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", token)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)
	var request *CreateRequest
	json.Unmarshal(respByte, &request)
	status = false
	orderid = ""
	if request.Status == 200 {
		status = true
		orderid = request.Data.Order.ID
	}
	return
}

func GetToken(url string) (token string, err error) {
	URL := strings.Join([]string{url, "GetColaToken"}, "/")
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
	var request *GetApiRequest
	json.Unmarshal(respByte, &request)
	if request.Status == 1 {
		return "", errors.New(request.Message)
	}
	return request.Token, nil
}

func AddAccount(url, oid string) bool {
	URL := strings.Join([]string{url, "AddAccount?account=", oid}, "/")
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
	var request *GetApiRequest
	json.Unmarshal(respByte, &request)
	if request.Status == 1 {
		return false
	}
	return true
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
		return "", err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}
	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	// decode image
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, nil)
	if err != nil {
		return "", err
	}
	// fmt.Println(result.String())
	return result.String(), err
}
