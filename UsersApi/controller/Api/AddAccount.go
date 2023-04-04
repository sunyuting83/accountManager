package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddAccount(c *gin.Context) {
	var (
		Account  string = c.Query("account")
		Password string = c.Query("password")
		status   string = c.DefaultQuery("status", "3")
		remarks  string = c.Query("remarks")
		IsJson   string = c.DefaultQuery("json", "0")
	)
	projectsID, ColaAPI := GetProjectsID(c)
	projectsInt, _ := strconv.Atoi(projectsID)
	NewStatus, _ := strconv.Atoi(status)
	if len(Account) == 0 {
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

	account := &database.Accounts{
		ProjectsID: uint(projectsInt),
		UserName:   Account,
		Password:   Password,
		NewStatus:  NewStatus,
	}
	if ColaAPI {
		account = &database.Accounts{
			ProjectsID: uint(projectsInt),
			UserName:   Account,
			NewStatus:  NewStatus,
			Remarks:    remarks,
		}
	}
	account.AddAccount()
	if IsJson == "1" {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "adding account successfully",
		})
		return
	}
	c.String(200, "添加成功")
}
