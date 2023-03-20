package main

import (
	"bytes"
	"colaAPI/ImageServer/Client/utils"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

type Status struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func postFile(filename, imgName, IMGServer string) bool {

	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return false
	}
	_, err = io.Copy(part, file)
	_ = writer.WriteField("name", imgName)
	_ = writer.WriteField("sk", "jLowHB51hBam4BhVBEduFye0WS1tjcCS")
	err = writer.Close()
	if err != nil {
		return false
	}

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	URL := strings.Join([]string{IMGServer, "set"}, "/")
	req, _ := http.NewRequest("POST", URL, body)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode > 200 {
		return false
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)

	var m *Status
	json.Unmarshal(respByte, &m)
	status := false
	if m.Status == 0 {
		status = true
	}
	return status
}

func postData(account, gold, multiple, diamond, crazy, cold, precise, exptime, imgName, APIServer string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("account", account)
	_ = writer.WriteField("gold", gold)
	_ = writer.WriteField("multiple", multiple)
	_ = writer.WriteField("diamond", diamond)
	_ = writer.WriteField("crazy", crazy)
	_ = writer.WriteField("cold", cold)
	_ = writer.WriteField("precise", precise)
	_ = writer.WriteField("exptime", exptime)
	_ = writer.WriteField("cover", imgName)
	err := writer.Close()
	if err != nil {
		fmt.Println("0")
		return
	}

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	URL := strings.Join([]string{APIServer, "PostSetAccount"}, "/")
	req, _ := http.NewRequest("POST", URL, body)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("0")
		return
	}
	if resp.StatusCode > 200 {
		fmt.Println("0")
		return
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)

	var m *Status
	json.Unmarshal(respByte, &m)
	status := "0"
	if m.Status == 0 {
		status = "1"
	}
	fmt.Println(status)
}

// sample usage
func main() {
	var (
		file     string
		account  string
		gold     string
		multiple string
		exptime  string
		diamond  string
		crazy    string
		cold     string
		precise  string
	)
	flag.StringVar(&file, "file", "3000", "image path")
	flag.StringVar(&account, "account", "0", "account")
	flag.StringVar(&gold, "gold", "0", "gold")
	flag.StringVar(&multiple, "multiple", "0", "multiple")
	flag.StringVar(&diamond, "diamond", "0", "Diamond")
	flag.StringVar(&crazy, "crazy", "0", "crazy")
	flag.StringVar(&cold, "cold", "0", "cold")
	flag.StringVar(&precise, "precise", "0", "precise")
	flag.StringVar(&exptime, "exptime", "0", "exptime")
	flag.Parse()

	OS := runtime.GOOS
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	CurrentPath, _ := utils.GetCurrentPath()
	fileName := strings.Join([]string{CurrentPath, file}, LinkPathStr)
	confYaml, err := utils.CheckConfig(OS, CurrentPath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	Key := GetProjectsKey(confYaml.APIServer)
	imgName := strings.Join([]string{Key, account}, "_")
	imgName = strings.Join([]string{imgName, "jpg"}, ".")
	status := postFile(fileName, imgName, confYaml.IMGServer)
	if status {
		postData(account, gold, multiple, diamond, crazy, cold, precise, exptime, imgName, confYaml.APIServer)
	}
}

func GetProjectsKey(url string) string {
	URList := strings.Split(url, "/")
	URLen := len(URList)
	key := URList[URLen-1]
	return key
}
