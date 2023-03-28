package controller

import (
	"colaAPI/database"
	"colaAPI/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AccountDrawedDateList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var date string = c.Query("date")
	var Limit string = c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)

	if len(date) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "错误的日期格式",
		})
		return
	}

	projectsID := GetProjectsID(c)
	ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)
	Projects, err := database.ProjectsCheckID(ProjectsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	startTime, endTime := utils.GetSqlDateTime(date)
	var acc *database.Accounts
	count, err := acc.GetDatedInCount(projectsID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	dataList, err := database.GetDatedInData(projectsID, startTime, endTime, pageInt, LimitInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if len(dataList) == 0 {
		dataList = make([]*database.Accounts, 0)
	}

	Data := gin.H{
		"status":   0,
		"data":     dataList,
		"projects": Projects,
		"total":    count,
	}
	c.JSON(http.StatusOK, Data)
}
