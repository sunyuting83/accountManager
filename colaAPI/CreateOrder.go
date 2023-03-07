package colaapi

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"
)

type RequestOrderList struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		List  []OrderList `json:"list"`
		Count int         `json:"count"`
	} `json:"data"`
}

type OrderList struct {
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

func CheckCreateOrder(c int, qrcode string) (newOrder bool) {
	newOrder = false
	logined, count := GetOrderList()
	if logined && count < c {
		newOrder = true
	}
	return
}

func GetOrderList() (logined bool, count int) {
	u := "Order/OrderList?type=5&state=0"
	url := strings.Join([]string{public.RootUrl, u}, "")
	count = 0
	platForm = &Platform{
		Method: "GET",
		URL:    url,
	}
	data, err := OrderListHttp(platForm)
	if err != nil {
		logined = false
	}
	if data.Status == "200" {
		logined = true
		count = data.Data.Count
	} else {
		logined = false
	}
	return
}

func OrderListHttp(post *Platform) (request RequestOrderList, err error) {
	req, _ := http.NewRequest(post.Method, post.URL, nil)
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
