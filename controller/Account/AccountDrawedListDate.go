package controller

import (
	"colaAPI/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllDateForAccountDrawed(c *gin.Context) {
	projectsID := GetProjectsID(c)
	ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)
	_, err := database.ProjectsCheckID(ProjectsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	dateList, err := database.GetDateTimeDataDraw(projectsID, "1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if len(dateList) == 0 {
		dateList = make([]string, 0)
	}

	Data := gin.H{
		"status":   0,
		"dateList": dateList,
	}
	c.JSON(http.StatusOK, Data)
}
