package controller

import (
	"colaAPI/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AccountList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var projectsID string = c.DefaultQuery("projectsID", "0")
	var Limit string = c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	var account *database.Accounts
	count, err := account.GetCount(projectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	ProjectsID, _ := strconv.Atoi(projectsID)
	Projects, err := database.ProjectsCheckID(int64(ProjectsID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	dataList, err := database.GetAccountList(pageInt, LimitInt, projectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	Data := gin.H{
		"status":   0,
		"data":     dataList,
		"total":    count,
		"projects": Projects,
	}
	c.JSON(http.StatusOK, Data)
}
