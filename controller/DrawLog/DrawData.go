package controller

import (
	"colaAPI/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DrawData(c *gin.Context) {
	var id string = c.Query("id")
	if len(id) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	data, err := database.GetDrawData(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	var account *[]database.Accounts
	json.Unmarshal([]byte(data.Data), &account)
	Data := gin.H{
		"status":   0,
		"data":     account,
		"projects": data.Projects,
	}
	c.JSON(http.StatusOK, Data)
}
