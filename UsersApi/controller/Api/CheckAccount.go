package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAccount(c *gin.Context) {
	var (
		account string = c.Query("account")
		IsJson  string = c.DefaultQuery("json", "0")
	)

	if len(account) == 0 {
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

	projectsID, _ := GetProjectsID(c)

	Projects, err := database.ProjectsCheckID(projectsID)
	if err != nil {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "参数错误",
				"project": Projects,
			})
		}
		c.String(200, "参数错误")
		return
	}

	Account, err := database.CheckOneAccount(projectsID, account)
	if err != nil {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "帐号不存在",
			})
		}
		c.String(200, "帐号不存在")
		return
	}
	if IsJson == "1" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "帐号存在",
			"data":    Account,
		})
		return
	}
	c.String(200, "帐号存在")
}
