package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetOrdersList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "50")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	UsersID := utils.GetCurrentUserID(c)
	user, err := database.UserCheckID(int64(UsersID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "get users failed",
		})
		return
	}
	var orders *database.Order
	count, err := orders.GetOrdersCount(user.ID)
	if count == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "还没有订单",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project",
		})
		return
	}
	dataList, err := database.GetOrdersList(pageInt, LimitInt, user.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project list",
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
