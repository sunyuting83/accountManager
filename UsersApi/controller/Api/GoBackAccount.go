package controller

import (
	BadgerDB "colaAPI/UsersApi/badger"
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GoBackAccount(c *gin.Context) {
	var (
		status  string = c.Query("status")
		IsJson  string = c.DefaultQuery("json", "0")
		windows string = c.DefaultQuery("windows", "0")
	)
	Path := c.Request.URL.Path
	PathList := strings.Split(Path, "/")
	Path = PathList[len(PathList)-1]
	if strings.Contains(Path, "clean") {
		status = GetBackPath(Path)
	}
	if !strings.Contains(Path, "clean") && len(status) == 0 {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "haven't params found",
			})
			return
		}
		c.String(200, "参数错误")
		return
	}
	var person Person
	c.ShouldBindUri(&person)
	Key := person.Key
	getnumber, _ := BadgerDB.Get([]byte(Key + ".getnumber"))
	getnumberInt, _ := strconv.Atoi(getnumber)
	if getnumberInt >= 300 {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "block",
			})
			return
		}
		c.String(200, "没有了")
		return
	}

	projectsID, _ := GetProjectsID(c)

	Projects, err := database.ProjectsCheckID(projectsID)
	if err != nil {
		if IsJson == "1" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		c.String(200, "参数错误")
		return
	}

	var statusJson []*StatusJSON
	json.Unmarshal([]byte(Projects.StatusJSON), &statusJson)

	var (
		hasPower     bool = false
		backToStatus string
	)

	for _, item := range statusJson {
		if item.Status == status {
			if item.CallBack {
				hasPower = true
				backToStatus = item.BackTo
				break
			}
		}
	}

	if !hasPower {
		if IsJson == "1" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": "状态不支持退回",
			})
			return
		}
		c.String(200, "状态不支持退回")
		return
	}
	backToStatusInt, _ := strconv.Atoi(backToStatus)
	var account *database.Accounts
	account.BackTo(projectsID, status, backToStatusInt, windows)
	if IsJson == "1" {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "退回成功",
		})
		return
	}
	c.String(200, "退回成功")
}
