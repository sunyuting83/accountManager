package controller

import (
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFiledList(c *gin.Context) {

	projectsID := GetProjectsID(c)

	Projects, err := database.ProjectsCheckID(projectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project ID",
		})
		return
	}
	filedList, err := database.GetFiledList(projectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't filed list",
		})
		return
	}

	var filedData string = ""
	if len(filedList) != 0 {
		filed, _ := database.GetFirstFiled(projectsID)
		filedData = filed.Data
	}
	var Sjson []*database.Accounts = make([]*database.Accounts, 0)
	json.Unmarshal([]byte(filedData), &Sjson)

	Data := gin.H{
		"status":   0,
		"dateList": filedList,
		"data":     Sjson,
		"projects": Projects,
	}
	c.JSON(http.StatusOK, Data)
}

func GetOneFiled(c *gin.Context) {
	var date string = c.Query("date")
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
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project ID",
		})
		return
	}

	filedData, _ := database.GetOneFiled(projectsID, date)
	var Sjson []*database.Accounts = make([]*database.Accounts, 0)
	json.Unmarshal([]byte(filedData.Data), &Sjson)

	Data := gin.H{
		"status":   0,
		"data":     Sjson,
		"projects": Projects,
	}
	c.JSON(http.StatusOK, Data)
}
