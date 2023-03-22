package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetAccount(c *gin.Context) {
	var To string = c.DefaultQuery("to", "1")
	var IsJson string = c.DefaultQuery("json", "0")
	var gameid string = c.Query("gameid")
	var Account string = c.Query("account")
	if len(Account) == 0 {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "Account must be specified",
			})
			return
		}
		c.String(200, "帐号不能为空")
		return
	}
	Path := c.Request.URL.Path
	PathList := strings.Split(Path, "/")
	Path = PathList[len(PathList)-1]
	if strings.Contains(Path, "finished") {
		To = GetSavePath(Path)
	}
	if gameid == "14e7110dd307" {
		To = "10"
	}

	var person Person
	var result *CacheValue
	if err := c.ShouldBindUri(&person); err != nil {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		c.String(200, "项目被锁定")
		return
	}
	var (
		projectsID string
	)
	has := Redis.Get(person.Key)
	if len(has) != 0 {
		json.Unmarshal([]byte(has), &result)
		projectsID = result.ProjectsID
	}
	account, err := database.CheckAccount(projectsID, Account)
	if err != nil {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		c.String(200, "帐号不存在")
		return
	}
	account.AccountUpStatus(To)
	if IsJson == "1" {
		Data := gin.H{
			"status":  0,
			"message": "成功",
		}
		c.JSON(http.StatusOK, Data)
		return
	}
	c.String(200, "成功")
}

func GetSavePath(s string) (to string) {
	switch s {
	case "regonefinished":
		to = "2"
	case "playonefinished":
		to = "4"
	case "billionfinished":
		to = "7"
	case "banfinished":
		to = "5"
	}
	return
}
