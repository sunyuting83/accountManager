package controller

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PostAccount struct {
	Account  string `form:"account" json:"account" xml:"account" binding:"required"`
	Gold     string `form:"gold" json:"gold" xml:"gold"  binding:"required"`
	Multiple string `form:"multiple" json:"multiple" xml:"multiple"`
	Diamond  string `form:"diamond" json:"diamond" xml:"diamond"`
	Crazy    string `form:"crazy" json:"crazy" xml:"crazy"`
	Cold     string `form:"cold" json:"cold" xml:"cold"`
	Precise  string `form:"precise" json:"precise" xml:"precise"`
	ExpTime  string `form:"exptime" json:"exptime" xml:"exptime"`
}

func PostSetAccount(c *gin.Context) {

	var form PostAccount
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if len(form.Gold) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  0,
			"message": "haven't node",
		})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "1上传文件失败",
		})
		return
	}

	b, _ := io.ReadAll(file)

	var gold int64
	if strings.Contains(form.Gold, "亿") {
		gx := strings.Split(form.Gold, "亿")
		goldstr := gx[0]
		if strings.Contains(form.Gold, ".") {
			g := strings.Split(goldstr, ".")
			// fmt.Println(g)
			goldstr = strings.Join([]string{g[0], g[1]}, "")
			// fmt.Println(goldstr)
			var x int64 = 10000000
			if len(g[1]) >= 2 {
				x = 1000000
			}
			n, _ := strconv.ParseInt(goldstr, 10, 64)
			// fmt.Println(n)
			// s := n * x
			gold = n * x
			// fmt.Println("yi1")
			// fmt.Println(gold)
		} else {
			// fmt.Println(goldstr)
			n, _ := strconv.ParseInt(goldstr, 10, 64)
			gold = n * 100000000
			// fmt.Println("yi2")
			// fmt.Println(gold)
		}
	} else if strings.Contains(form.Gold, "万") {
		gx := strings.Split(form.Gold, "万")
		goldstr := gx[0]
		n, _ := strconv.ParseInt(goldstr, 10, 64)
		gold = n * 10000
		// fmt.Println("wan2")
		// fmt.Println(gold)
	} else {
		n, _ := strconv.ParseInt(form.Gold, 10, 64)
		gold = n
	}

	c.String(200, "成功")
}

func postFile(file multipart.File) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", file)

	io.Copy(part, file)
	_ = writer.WriteField("name", gold)
	writer.Close()

	client := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}
	req, _ := http.NewRequest("POST", "http://localhost:13005/set", body)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	// return nil
}
