package controller

import (
	"colaAPI/database"
	"colaAPI/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UsersList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	var user database.Users
	count, err := user.GetCount()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	result := utils.GetTokenUserData(c)
	dataList, err := database.GetUsersList(pageInt, LimitInt, result.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	Data := gin.H{
		"status": 0,
		"data":   dataList,
		"total":  count,
	}
	c.JSON(http.StatusOK, Data)
}
func UsersAllList(c *gin.Context) {
	result := utils.GetTokenUserData(c)
	dataList, err := database.GetAllUsersList(result.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	Data := gin.H{
		"status": 0,
		"data":   dataList,
	}
	c.JSON(http.StatusOK, Data)
}
