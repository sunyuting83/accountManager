package controller

import (
	BadgerDB "colaAPI/Users/badger"
	"colaAPI/Users/database"
	geo "colaAPI/Users/geoip"
	"colaAPI/Users/utils"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type RegeditForm struct {
	UserName   string `form:"username" json:"username" xml:"username"  binding:"required"`
	Password   string `form:"password" json:"password" xml:"password"  binding:"required"`
	RePassword string `form:"repassword" json:"repassword" xml:"repassword"  binding:"required"`
	VCode      string `form:"vcode" json:"vcode" xml:"vcode"  binding:"required"`
}

func Regedit(c *gin.Context) {
	var form RegeditForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if len(form.UserName) < 4 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't username",
		})
		return
	}
	if utils.ContainsSpecialCharacters(form.UserName) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "用户名只可以是英文加数字",
		})
		return
	}
	ignoreUserName := "admin|manage"
	if strings.ContainsAny(form.UserName, ignoreUserName) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "非法用户名",
		})
		return
	}
	if len(form.Password) < 8 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't password",
		})
		return
	}
	if len(form.VCode) < 5 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't verify code",
		})
		return
	}
	if form.Password != form.RePassword {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "password error",
		})
		return
	}
	VCode := utils.ConvertToUpperCase(form.VCode)

	vcode, err := BadgerDB.Get([]byte(VCode))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "验证码错误",
		})
		return
	}
	if vcode != VCode {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "验证码错误",
		})
		return
	}
	ipStr := c.ClientIP()
	ip := net.ParseIP(ipStr)

	IPCount, err := BadgerDB.Get([]byte(ipStr))
	if err != nil && err.Error() != "Key not found" {
		var ttl int64 = 60 * 60 * 24 // ttl以秒为单位
		BadgerDB.SetWithTTL([]byte(ipStr), []byte("1"), ttl)
		IPCount = "1"
	}
	ipCount, _ := strconv.Atoi(IPCount)
	if ipCount != 0 && ipCount >= 5 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "同IP在1天内只能注册5个帐号",
		})
		return
	}
	// 获取IP 写入Cache Key IP地址 Value 已注册数量 如果数量超过5返回错误 1天内同IP只能注册5个帐号 超时 1天

	user, err := database.CheckUserName(form.UserName)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if len(user.UserName) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "用户名已存在",
		})
		return
	}
	secret_key, _ := c.Get("secret_key")
	SECRET_KEY := secret_key.(string)
	PASSWD := utils.MD5(strings.Join([]string{form.Password, SECRET_KEY}, ""))

	LocalAddress := "未知"

	if ipStr == "127.0.0.1" {
		LocalAddress = "本机"
	}

	if ipStr != "127.0.0.1" {
		geoDB, err := geo.LoadGeoFile()
		if err != nil {
			LocalAddress = "未知"
		}
		defer geoDB.Close()
		record, err := geoDB.City(ip)
		if err != nil {
			LocalAddress = "未知"
		}

		// 读取位置信息
		cityName := record.City.Names["cn"]
		countryName := record.Country.Names["cn"]
		LocalAddress = strings.Join([]string{countryName, cityName}, "")
	}

	tiemNow := time.Now().Format("2006-01-02_15:04:05")
	WalletAddress := utils.MD5(strings.Join([]string{ipStr, form.UserName, tiemNow, form.Password, SECRET_KEY}, ""))
	WalletAddress = WalletAddress[:16]

	var adduser *database.CoinUsers
	adduser.UserName = form.UserName
	adduser.Password = PASSWD
	adduser.NewStatus = 0
	adduser.Coin = 0.0
	adduser.IPAddress = ipStr
	adduser.LocalAddress = LocalAddress
	adduser.WalletAddress = WalletAddress
	err = adduser.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		BadgerDB.Delete([]byte(VCode))
		return
	}
	ipCount = ipCount + 1
	ipCountStr := strconv.Itoa(ipCount)
	BadgerDB.UpdateWithOutTTL([]byte(ipStr), []byte(ipCountStr))
	BadgerDB.Delete([]byte(VCode))
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "注册成功",
		"data":    user,
	})
}
