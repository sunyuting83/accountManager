package colaapi

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
