package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SearchProducts(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "50")
	var gold string = c.DefaultQuery("gold", "0")
	var multiple string = c.DefaultQuery("multiple", "0")
	var diamond string = c.DefaultQuery("diamond", "0")
	var crazy string = c.DefaultQuery("crazy", "0")
	var cold string = c.DefaultQuery("cold", "0")
	var precise string = c.DefaultQuery("precise", "0")
	var gameid string = c.Query("gameid")
	if len(gameid) <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "必须包含游戏ID",
		})
		return
	}
	if !isNumeric(gameid) {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "非法请求",
		})
		return
	}
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	GameID64, _ := strconv.Atoi(gameid)
	GameID := uint(GameID64)
	Gold, _ := strconv.ParseInt(gold, 10, 64)
	Multiple, _ := strconv.ParseInt(multiple, 10, 64)
	Diamond, _ := strconv.ParseInt(diamond, 10, 64)
	Crazy, _ := strconv.ParseInt(crazy, 10, 64)
	Cold, _ := strconv.ParseInt(cold, 10, 64)
	Precise, _ := strconv.ParseInt(precise, 10, 64)

	form := &utils.Filter{
		MinGold:  Gold,
		Multiple: Multiple,
		Diamond:  Diamond,
		Crazy:    Crazy,
		Cold:     Cold,
		Precise:  Precise,
	}
	var account *database.Accounts
	count, err := account.GetCountUseScopesB(form, pageInt, LimitInt, GameID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project",
		})
		return
	}
	rows, err := database.GetDataUseScopesB(form, pageInt, LimitInt, GameID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project",
		})
		return
	}

	DataList := MakeDataList(rows)
	Data := gin.H{
		"status": 0,
		"data":   DataList,
		"total":  count,
	}
	c.JSON(http.StatusOK, Data)
}
func isNumeric(str string) bool {
	_, err := strconv.ParseFloat(str, 64)
	return err == nil
}
