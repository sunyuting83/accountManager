package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	result := utils.GetTokenUserData(c)

	user, err := database.UserCheckID(int64(result.UserID))
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	user.Password = ""

	Data := gin.H{
		"status":  0,
		"message": "获取成功",
		"users":   user,
	}
	c.JSON(http.StatusOK, Data)
}
