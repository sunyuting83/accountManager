package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddAccount(c *gin.Context) {
	var Account string = c.Query("account")
	var status string = c.DefaultQuery("status", "3")
	var remarks string = c.Query("remarks")
	projectsID, ColaAPI := GetProjectsID(c)
	projectsInt, _ := strconv.Atoi(projectsID)
	NewStatus, _ := strconv.Atoi(status)
	if ColaAPI {
		account := &database.Accounts{
			ProjectsID: uint(projectsInt),
			UserName:   Account,
			NewStatus:  NewStatus,
			Remarks:    remarks,
		}
		account.AddAccount()
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "adding account successfully",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}
