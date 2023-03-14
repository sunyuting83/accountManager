package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectsModify struct {
	UserName  string `form:"username" json:"username" xml:"username"`
	Password  string `form:"password" json:"password" xml:"password"`
	AccNumber int    `form:"AccNumber" json:"AccNumber" xml:"AccNumber"`
}

func ModifyProjects(c *gin.Context) {
	var form ProjectsModify
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	var person Person
	var result *CacheValue
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	var (
		ProjectsID string
	)
	has := Redis.Get(person.Key)
	if len(has) != 0 {
		json.Unmarshal([]byte(has), &result)
		ProjectsID = result.ProjectsID
	}
	projects, err := database.ProjectsCheckID(ProjectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project ID",
		})
		return
	}
	if !projects.ColaAPI {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "it's not a cola api projects",
		})
		return
	}

	if len(form.UserName) < 5 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't username for API",
		})
		return
	}
	if len(form.Password) < 6 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't Password for API",
		})
		return
	}
	if form.AccNumber == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't AccNumber for API",
		})
		return
	}
	projects = &database.Projects{
		UserName:  form.UserName,
		Password:  form.Password,
		AccNumber: form.AccNumber,
	}
	projects.UpdateProjects(ProjectsID)
	projects, err = database.ProjectsCheckID(ProjectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "修改成功",
		"data":    projects,
	})
}
