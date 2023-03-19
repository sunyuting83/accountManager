package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type Publication struct {
	RootUrl string
}

type Platform struct {
	Method string
	URL    string
	Params string
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

type CreateData struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   Data   `json:"data"`
}
type List struct {
	OID           string `json:"o_id"`
	OrderID       string `json:"order_id"`
	UserID        string `json:"user_id"`
	AppName       string `json:"app_name"`
	Wxid          string `json:"wxid"`
	Amount        string `json:"amount"`
	Time          string `json:"time"`
	State         string `json:"state"`
	Type          string `json:"type"`
	PType         string `json:"p_type"`
	Projectid     string `json:"projectid"`
	NewTime       string `json:"new_time"`
	WxNickname    string `json:"wx_nickname"`
	Remarks       string `json:"remarks"`
	Queta         string `json:"queta"`
	JType         string `json:"j_type"`
	Other         string `json:"other"`
	ChannelID     string `json:"channel_id"`
	Apigateway    string `json:"apigateway"`
	Remoteorderid string `json:"remoteorderid"`
	AppLogo       string `json:"app_logo"`
}
type Data struct {
	List  []List `json:"list"`
	Count int    `json:"count"`
}

var public *Publication = &Publication{
	RootUrl: "http://tiancaiapi.tablecando.cn/api/",
}
var platForm *Platform

func Login(username string, password string) (token string, err error) {
	u := "Member/login"
	url := strings.Join([]string{public.RootUrl, u}, "")

	Params := strings.Join([]string{"member_user=", username, "&member_pwd=", password}, "")
	platForm = &Platform{
		Method: "POST",
		URL:    url,
		Params: Params,
	}
	data, err := LoginHttp(platForm)
	if err != nil {
		return "", err
	}
	if data.Status == "200" {
		return data.Token, nil
	} else {
		return "", nil
	}
}

func LoginHttp(post *Platform) (request RequestLogin, err error) {
	req, _ := http.NewRequest(post.Method, post.URL, strings.NewReader(post.Params))
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
	return request, nil
}

func GetOrderListNumber(token string) (n int, status string, err error) {
	URL := "http://tiancaiapi.tablecando.cn/api/Order/OrderList?type=5&state=0"
	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Authorization", token)
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	resp, err := (&http.Client{Timeout: 35 * time.Second}).Do(req)
	if err != nil {
		return 0, "0", err
	}
	defer resp.Body.Close()
	respByte, _ := io.ReadAll(resp.Body)
	var request *CreateData
	json.Unmarshal(respByte, &request)
	if request.Status == "200" {
		return request.Data.Count, request.Status, nil
	}
	return
}
