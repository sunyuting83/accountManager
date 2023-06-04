package controller

import (
	"colaAPI/UsersApi/database"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
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

type SqlData struct {
	Cover         string  `json:"cover" structs:"cover"`
	TodayGold     int64   `json:"today_gold" structs:"today_gold"`
	Multiple      int64   `json:"multiple" structs:"multiple"`
	Diamond       int     `json:"diamond" structs:"diamond"`
	Crazy         int     `json:"crazy" structs:"crazy"`
	Cold          int     `json:"cold" structs:"cold"`
	Precise       int     `json:"precise" structs:"precise"`
	Exptime       int64   `json:"exptime" structs:"exptime"`
	Price         float64 `json:"price" structs:"price"`
	YesterdayGold int64   `json:"yesterday_gold" structs:"yesterday_gold"`
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
			"status":  1,
			"message": "haven't gold",
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

	account, err := database.CheckOneAccount(projectsID, form.Account)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "帐号不存在",
		})
		return
	}

	project, err := database.FindGames(projectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "帐号不存在",
		})
		return
	}
	// fmt.Println(project.Games.UnitPrice / float64(project.Games.SingleNumber*100000000))
	Price := Decimal(project.Games.BasePrice + ((project.Games.UnitPrice / float64(project.Games.SingleNumber*100000000)) * float64(gold)))
	// fmt.Println(project.Games.BasePrice, project.Games.UnitPrice, project.Games.SingleNumber, float64(gold))
	// fmt.Println(Price)
	updata := &SqlData{
		Cover:     form.Cover,
		TodayGold: gold,
		Multiple:  Multiple,
		Diamond:   Diamond,
		Crazy:     Crazy,
		Cold:      Cold,
		Precise:   Precise,
		Exptime:   ExpTimeInt,
		Price:     Price,
	}

	timeobj := time.Unix(account.UpdatedAt/1000, 0)
	olDate := timeobj.Format("20060102")
	nowTime := time.Now()
	timeStr := nowTime.Format("20060102")
	// fmt.Println(olDate, timeStr)
	if timeStr > olDate {
		updata.YesterdayGold = account.TodayGold
	}

	updataMAP := structs.Map(&updata)
	account.UpdataOneAccount(projectsID, form.Account, updataMAP)

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
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

func Decimal(num float64) float64 {
	num, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", num), 64)
	return num
}
