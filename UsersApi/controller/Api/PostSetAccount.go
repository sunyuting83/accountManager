package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PostAccount struct {
	Account  string `form:"account" json:"account" xml:"account"  binding:"required"`
	Gold     string `form:"gold" json:"gold" xml:"gold"  binding:"required"`
	Multiple string `form:"multiple" json:"multiple" xml:"multiple"`
	Cover    string `form:"cover" json:"cover" xml:"cover"`
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

	Multiple, _ := strconv.ParseInt(form.Multiple, 10, 64)
	Diamond, _ := strconv.Atoi(form.Diamond)
	Crazy, _ := strconv.Atoi(form.Crazy)
	Cold, _ := strconv.Atoi(form.Cold)
	Precise, _ := strconv.Atoi(form.Precise)

	ExpTimeInt := strToDate(form.ExpTime)
	projectsID, _ := GetProjectsID(c)

	account, err := database.CheckAccount(projectsID, form.Account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "帐号不存在",
		})
		return
	}
	updata := database.Accounts{
		Cover:     form.Cover,
		TodayGold: gold,
		Multiple:  Multiple,
		Diamond:   Diamond,
		Crazy:     Crazy,
		Cold:      Cold,
		Precise:   Precise,
		Exptime:   ExpTimeInt,
	}

	timeobj := time.Unix(int64(account.UpdatedAt), 0)
	olDate := timeobj.Format("20060102")
	nowTime := time.Now()
	timeStr := nowTime.Format("20060102")
	if timeStr > olDate {
		updata.YesterdayGold = account.TodayGold
	}

	account.UpdataOneAccount(projectsID, form.Account, updata)

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "上传文件成功",
	})
}
func strToDate(date string) (d int64) {
	var LOC, _ = time.LoadLocation("Asia/Shanghai")
	res1, err := time.ParseInLocation("2006/01/02 15:04:05", date, LOC)
	if err != nil {
		return 0
	}
	return res1.Unix()
}
