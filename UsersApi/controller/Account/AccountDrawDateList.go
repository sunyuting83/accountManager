package controller

import (
	"colaAPI/UsersApi/database"
	"colaAPI/UsersApi/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AccountDrawDateList(c *gin.Context) {
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

	Projects, err := database.ProjectsCheckID(projectsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	var statusJson []*StatusJSON
	json.Unmarshal([]byte(Projects.StatusJSON), &statusJson)

	var (
		hasStatus []string
	)
	for _, item := range statusJson {
		if item.Pull {
			hasStatus = append(hasStatus, item.Status)
		}
	}

	startTime, endTime := utils.GetSqlDateTime(date)
	var acc *database.Accounts
	count, err := acc.GetDateInCount(projectsID, hasStatus, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	dataList, err := database.GetDateInData(projectsID, hasStatus, startTime, endTime, pageInt, LimitInt)
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
	games, _ := database.GetFirstCalc(Projects.GamesID)

	Data := gin.H{
		"status":   0,
		"data":     dataList,
		"projects": Projects,
		"total":    count,
		"games":    games,
	}
	c.JSON(http.StatusOK, Data)
}
