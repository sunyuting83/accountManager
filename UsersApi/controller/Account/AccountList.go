package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Key string `uri:"key" binding:"required"`
}
type CacheValue struct {
	UsersID    string `json:"UsersID"`
	ProjectsID string `json:"ProjectsID"`
}

func AccountList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Status string = c.DefaultQuery("status", "0")
	var Limit string = c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)

	projectsID := GetProjectsID(c)

	var account *database.Accounts
	count, err := account.GetCount(projectsID, Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project",
		})
		return
	}
	Projects, err := database.ProjectsCheckID(projectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project ID",
		})
		return
	}
	dataList, err := database.GetAccountList(pageInt, LimitInt, projectsID, Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project list",
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

func GetProjectsID(c *gin.Context) (projectsID string) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	var result *CacheValue
	has := Redis.Get(person.Key)
	if len(has) != 0 {
		json.Unmarshal([]byte(has), &result)
		projectsID = result.ProjectsID
	}
	return
}
