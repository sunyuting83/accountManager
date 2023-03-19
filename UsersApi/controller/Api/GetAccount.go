package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/UsersApi/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Key string `uri:"key" binding:"required"`
}
type CacheValue struct {
	UsersID    string `json:"UsersID"`
	ProjectsID string `json:"ProjectsID"`
	ColaAPI    bool   `json:"ColaAPI"`
}

func GetOneAccount(c *gin.Context) {
	var (
		computid string = c.Query("computid")
		status   string = c.Query("status")
		to       string = c.Query("to")
		IsJson   string = c.DefaultQuery("json", "0")
		splitStr string = c.DefaultQuery("splitStr", "----")
	)
	Path := c.Request.URL.Path
	PathList := strings.Split(Path, "/")
	Path = PathList[len(PathList)-1]
	if len(Path) == 6 {
		status, to = GetPath(Path)
	}
	if len(status) == 0 {
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
	if len(to) == 0 {
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
	projectsID, ColaAPI := GetProjectsID(c)
	var (
		comput *database.Comput
		err    error
	)
	fmt.Println(ColaAPI)
	if len(computid) > 28 {
		comput, err = database.GetOneComputer(computid)
		if err != nil {
			comput = &database.Comput{
				ComputCode: computid,
			}
			comput.ComputerInsert()
		}
	}

	account, err := database.GetOneAccount(projectsID, status)
	if err != nil {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "haven't project list",
			})
			return
		}
		c.String(200, "没有了")
		return
	}
	account.AccountUpStatus(to)
	if len(computid) > 28 {
		account.AccountUpComput(comput.ID)
	}

	if IsJson == "1" {
		Data := gin.H{
			"status": 0,
			"data":   account,
		}
		c.JSON(http.StatusOK, Data)
		return
	}
	AccountString := strings.Join([]string{account.UserName, account.Password}, splitStr)
	if len(account.PhoneNumber) != 0 {
		AccountString = strings.Join([]string{account.PhoneNumber}, splitStr)
	}
	if len(account.PhonePassword) != 0 {
		AccountString = strings.Join([]string{account.PhonePassword}, splitStr)
	}
	c.String(200, AccountString)
}

func GetProjectsID(c *gin.Context) (projectsID string, ColaAPI bool) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	var result *CacheValue
	has := Redis.Get(person.Key)
	if len(has) != 0 {
		json.Unmarshal([]byte(has), &result)
		projectsID = result.ProjectsID
		ColaAPI = result.ColaAPI
	}
	return
}

func GetPath(s string) (status, to string) {
	switch s {
	case "findregone":
		status = "0"
		to = "1"
	case "getplayone":
		status = "2"
		to = "3"
	}
	return
}
