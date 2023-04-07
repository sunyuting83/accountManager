package controller

import (
	controller "colaAPI/controller/Account"
	"colaAPI/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DrawList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	projectsID := controller.GetProjectsID(c)
	var draw database.DrawLogs
	count, err := draw.GetCount(projectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	dataList, err := database.GetDrawList(pageInt, LimitInt, projectsID)
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
