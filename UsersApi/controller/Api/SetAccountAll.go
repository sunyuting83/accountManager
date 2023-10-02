package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetAccountAll(c *gin.Context) {
	var (
		Status   string = c.Query("status")
		Password string = c.Query("password")
		Remarks  string = c.Query("remarks")
		Gold     string = c.Query("gold")
		Multiple string = c.Query("multiple")
		Cover    string = c.Query("cover")
		Diamond  string = c.Query("diamond")
		Crazy    string = c.Query("crazy")
		Cold     string = c.Query("cold")
		Precise  string = c.Query("precise")
		ExpTime  string = c.Query("exptime")
		IsJson   string = c.DefaultQuery("json", "0")
		Account  string = c.Query("account")
	)
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
	account, err := database.CheckOneAccount(projectsID, Account)
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
	if account.NewStatus == 108 {
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
	UpData := make(map[string]interface{}, 1)
	if len(Status) != 0 {
		StatusInt, _ := strconv.Atoi(Status)
		UpData["NewStatus"] = StatusInt
	}
	if len(Password) != 0 {
		UpData["Password"] = Password
	}
	if len(Remarks) != 0 {
		UpData["Remarks"] = Remarks
	}
	if len(Gold) != 0 {
		GoldInt, _ := strconv.ParseInt(Gold, 10, 64)
		UpData["TodayGold"] = GoldInt
	}
	if len(Multiple) != 0 {
		MultipleInt, _ := strconv.ParseInt(Multiple, 10, 64)
		UpData["Multiple"] = MultipleInt
	}
	if len(Cover) != 0 {
		UpData["Cover"] = Cover
	}
	if len(Diamond) != 0 {
		DiamondInt, _ := strconv.Atoi(Diamond)
		UpData["Diamond"] = DiamondInt
	}
	if len(Crazy) != 0 {
		CrazyInt, _ := strconv.Atoi(Crazy)
		UpData["Crazy"] = CrazyInt
	}
	if len(Cold) != 0 {
		ColdInt, _ := strconv.Atoi(Cold)
		UpData["Cold"] = ColdInt
	}
	if len(Precise) != 0 {
		PreciseInt, _ := strconv.Atoi(Precise)
		UpData["Precise"] = PreciseInt
	}
	if len(ExpTime) != 0 {
		ExpTimeInt := strToDate(ExpTime)
		UpData["Exptime"] = ExpTimeInt
	}

	account.AccountUpAll(UpData)
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
