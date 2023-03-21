package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddAccount(c *gin.Context) {
	var Account string = c.Query("account")
	projectsID, ColaAPI := GetProjectsID(c)
	projectsInt, _ := strconv.Atoi(projectsID)

	if ColaAPI {
		account := &database.Accounts{
			ProjectsID: uint(projectsInt),
			UserName:   Account,
			NewStatus:  3,
		}
		account.AddAccount()
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "adding account successfully",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}
