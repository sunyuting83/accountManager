package controller

import (
	"colaAPI/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	var admin database.Manager
	count, err := admin.GetCount()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	dataList, err := database.GetAdminList(pageInt, LimitInt)
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
