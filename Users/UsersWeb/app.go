package main

import (
	BadgerDB "UsersWeb/badger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

const RootURL = "http://localhost:13006/api/v1/"

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) AddCart(c []int) map[string]interface{} {
	var errResponse = make(map[string]interface{})
	errResponse["status"] = 1
	// BadgerDB.Delete([]byte("cart"))
	if len(c) == 0 {
		errResponse["message"] = "参数不能为空"
		return errResponse
	}
	newCartList := RemoveRepeatedList(c)
	cart, err := BadgerDB.GetToken([]byte("cart"))
	if err != nil {
		paramsByte, _ := json.Marshal(newCartList)
		BadgerDB.Set([]byte("cart"), paramsByte)
		errResponse["status"] = 0
		errResponse["message"] = "添加成功"
		return errResponse
	}
	var inttt []int
	json.Unmarshal(cart, &inttt)
	CartList := FilterUniqueValues(newCartList, inttt)
	// CartListByte := IntSliceToByteSlice(CartList)
	CartListByte, _ := json.Marshal(CartList)
	BadgerDB.Set([]byte("cart"), CartListByte)
	errResponse["status"] = 0
	errResponse["message"] = "添加成功"
	return errResponse
}

func (a *App) CleanCart() map[string]interface{} {
	var Response = make(map[string]interface{})
	Response["status"] = 0
	Response["message"] = "成功清空"
	Response["data"] = make([]interface{}, 0)
	// BadgerDB.Delete([]byte("cart"))
	// return errResponse
	BadgerDB.Delete([]byte("cart"))
	return Response
}

func (a *App) GetCart() map[string]interface{} {
	var errResponse = make(map[string]interface{})
	errResponse["status"] = 0
	// BadgerDB.Delete([]byte("cart"))
	// return errResponse
	cart, err := BadgerDB.GetToken([]byte("cart"))
	if err != nil {
		errResponse["data"] = make([]map[string]interface{}, 0)
		return errResponse
	}
	var inttt []int
	json.Unmarshal(cart, &inttt)
	params := make(map[string]interface{})
	params["cart"] = inttt
	data := HTTPRequest("POST", "GetCart", params)
	return data
}

func (a *App) DeleteCart(ID int) map[string]interface{} {
	var errResponse = make(map[string]interface{})
	errResponse["status"] = 1
	cart, err := BadgerDB.GetToken([]byte("cart"))
	if err != nil {
		errResponse["message"] = "发生错误"
		return errResponse
	}
	var inttt []int
	json.Unmarshal(cart, &inttt)
	newList := filterValue(inttt, ID)
	paramsByte, _ := json.Marshal(newList)
	BadgerDB.Set([]byte("cart"), paramsByte)
	errResponse["status"] = 0
	errResponse["message"] = "删除成功"
	return errResponse
}

func (a *App) PostOrders(ids []int) map[string]interface{} {
	var errResponse = make(map[string]interface{})
	errResponse["status"] = 0
	params := make(map[string]interface{})
	params["list"] = ids
	data := HTTPRequest("POST", "PostOrders", params)
	return data
}

func (a *App) GetOrdersList(params map[string]interface{}) map[string]interface{} {
	data := HTTPRequest("GET", "GetOrdersList", params)
	return data
}

func (a *App) GetOrdersDetail(params map[string]interface{}) map[string]interface{} {
	data := HTTPRequest("GET", "GetOrdersDetail", params)
	return data
}

func (a *App) OrderRefund(params map[string]interface{}) map[string]interface{} {
	data := HTTPRequest("POST", "OrderRefund", params)
	return data
}

func (a *App) FormatDateTime(timestamp int64) string {
	t := time.Unix(timestamp/1000, 0)
	return t.Format("2006-01-02 15:04:05")
}

func (a *App) GetGamesList() map[string]interface{} {
	// var Response = make(map[string]interface{})
	// Response["status"] = 1
	// BadgerDB.Delete([]byte("gamelist"))
	// return Response
	gamelist, err := BadgerDB.GetToken([]byte("gamelist"))
	if err == nil {
		map2 := make(map[string]interface{})
		json.Unmarshal(gamelist, &map2)
		return map2
	}
	params := make(map[string]interface{})
	data := HTTPRequest("GET", "GetGamesList", params)
	dataJSON, _ := json.Marshal(data)
	BadgerDB.SetWithTTL([]byte("gamelist"), dataJSON, 60*60*24*3)
	return data
}

func (a *App) GetProducts(params map[string]interface{}) map[string]interface{} {
	data := HTTPRequest("GET", "GetProducts", params)
	return data
}

func (a *App) CheckLogin() map[string]interface{} {
	var errResponse = make(map[string]interface{})
	errResponse["status"] = 1
	token, err := BadgerDB.Get([]byte("token"))
	if err != nil {
		errResponse["message"] = err.Error()
		return errResponse
	}
	errResponse["token"] = token
	params := make(map[string]interface{})
	data := HTTPRequest("GET", "CheckLogin", params)
	return data
}

func (a *App) Login(params map[string]interface{}) map[string]interface{} {
	data := HTTPRequest("POST", "Login", params)
	if Strval(data["status"]) == "0" {
		token := Strval(data["token"])
		BadgerDB.SetWithTTL([]byte("token"), []byte(token), 60*60*24*30)
	}
	return data
}

func (a *App) Logout() map[string]interface{} {
	BadgerDB.Delete([]byte("token"))
	var errResponse = make(map[string]interface{})
	errResponse["status"] = 0
	return errResponse
}

func (a *App) GetUsers() map[string]interface{} {
	params := make(map[string]interface{})
	data := HTTPRequest("GET", "GetUsers", params)
	return data
}

func HTTPRequest(method, uri string, params map[string]interface{}) map[string]interface{} {
	var req *http.Request
	var err error
	var token string
	var errResponse = make(map[string]interface{})
	errResponse["status"] = 1

	token, err = BadgerDB.Get([]byte("token"))
	if err != nil && err.Error() != "Key not found" {
		errResponse["message"] = "403"
		return errResponse
	}

	token = strings.Join([]string{"Bearer", token}, " ")

	uri = strings.Join([]string{RootURL, uri}, "")
	// 根据方法选择请求类型
	switch strings.ToUpper(method) {
	case "GET":
		// GET 请求
		getRequestParams := buildQueryParams(params)
		getURL := uri + getRequestParams
		req, _ = http.NewRequest("GET", getURL, nil)
	case "POST":
		// POST 请求
		paramsByte, _ := json.Marshal(params)
		data := bytes.NewReader(paramsByte)
		req, _ = http.NewRequest("POST", uri, data)
	case "PUT":
		// PUT 请求
		paramsByte, _ := json.Marshal(params)
		data := bytes.NewReader(paramsByte)
		req, _ = http.NewRequest("PUT", uri, data)
	case "DELETE":
		// DELETE 请求
		paramsByte, _ := json.Marshal(params)
		data := bytes.NewReader(paramsByte)
		req, _ = http.NewRequest("DELETE", uri, data)
	default:
		return errResponse
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", token)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	// 发送请求并返回响应
	client := http.Client{Timeout: 35 * time.Second}
	resp, err := client.Do(req)
	if err != nil && resp.StatusCode >= 400 {
		errResponse["message"] = "403"
		return errResponse
	}
	if err != nil {
		errResponse["message"] = err.Error()
		return errResponse
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		errResponse["message"] = err.Error()
		return errResponse
	}
	if len(respByte) == 0 {
		errResponse["message"] = "发生错误"
		return errResponse
	}

	return parseResponse(respByte)
}

// 解析响应为 map[string]interface{}
func parseResponse(respByte []byte) map[string]interface{} {
	respMap := make(map[string]interface{})
	err := json.Unmarshal(respByte, &respMap)
	if err != nil {
		return nil
	}
	return respMap
}

func buildQueryParams(params map[string]interface{}) string {
	var queryParams []string
	for key, value := range params {
		strParam := fmt.Sprintf("%s=%v", key, value)
		queryParams = append(queryParams, strParam)
	}

	joinedParams := strings.Join(queryParams, "&")
	return strings.Join([]string{"?", joinedParams}, "")
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

func RemoveRepeatedList(personList []int) (result []int) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < n; j++ {
			if personList[i] == personList[j] {
				repeat = true
				break
			}
		}
		if !repeat && personList[i] != 0 {
			result = append(result, personList[i])
		}
	}
	return
}

func FilterUniqueValues(arr1, arr2 []int) []int {
	uniqueValues := make(map[int]bool)
	for _, num := range arr1 {
		uniqueValues[num] = true
	}
	for _, num := range arr2 {
		uniqueValues[num] = true
	}
	result := make([]int, 0, len(uniqueValues))
	for num := range uniqueValues {
		result = append(result, num)
	}
	return result
}

func filterValue(arr []int, value int) []int {
	var result []int
	for _, v := range arr {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}
