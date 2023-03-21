package controller

import (
	BadgerDB "colaAPI/UsersApi/badger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetColaToken(c *gin.Context) {
	projectsID, ColaAPI := GetProjectsID(c)
	if ColaAPI {
		token, err := BadgerDB.Get([]byte(projectsID + ".token"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "get token succeeded",
			"token":   token,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}
