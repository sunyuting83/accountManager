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
	projectsID, _ := GetProjectsID(c)
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

	project, err := database.ProjectsCheckID(projectsID)
	if err != nil {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "database error",
			})
			return
		}
		c.String(200, "数据库错误")
		return
	}
	acc, err := database.CheckOneAccount(projectsID, Account)
	if err != nil && err.Error() != "record not found" {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "database error",
			})
			return
		}
		c.String(200, "数据库错误")
		return
	}
	if acc.ID > 0 {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "has account",
			})
			return
		}
		c.String(200, "帐号已存在")
		return
	}
	if err != nil && err.Error() == "record not found" {
		account := &database.Accounts{
			ProjectsID: uint(projectsInt),
			GameID:     &project.GamesID,
			UserName:   Account,
			Password:   Password,
			NewStatus:  NewStatus,
		}
		if len(Password) != 0 {
			account.Password = Password
		}
		if len(remarks) != 0 {
			account.Remarks = remarks
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
		return
	}
	if IsJson == "1" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "database error",
		})
		return
	}
	c.String(200, "数据库错误")
}
