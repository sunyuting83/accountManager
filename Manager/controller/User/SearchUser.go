package controller

import (
	"colaAPI/Manager/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchUser(c *gin.Context) {
	var key string = c.DefaultQuery("key", "0")

	data, err := database.SearchUsers(key)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	Data := gin.H{
		"status": 0,
		"data":   data,
	}
	c.JSON(http.StatusOK, Data)
}
