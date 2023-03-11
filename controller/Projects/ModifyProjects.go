package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectsModify struct {
	ID         string `form:"id" json:"id" xml:"id"  binding:"required"`
	UserName   string `form:"username" json:"username" xml:"username"`
	Password   string `form:"password" json:"password" xml:"password"`
	AccNumber  int    `form:"AccNumber" json:"AccNumber" xml:"AccNumber"`
	ColaAPI    string `form:"ColaAPI" json:"ColaAPI" xml:"ColaAPI"`
	StatusJSON string `form:"StatusJSON" json:"StatusJSON" xml:"StatusJSON"`
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

	StatusJSON := MakeStatusJSON(form.StatusJSON)

	var ColaAPI1 bool = false
	if form.ColaAPI == "true" {
		ColaAPI1 = true
	}
	if ColaAPI1 {
		if len(form.UserName) < 5 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  1,
				"message": "haven't projects name",
			})
			return
		}
		if len(form.Password) < 6 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  1,
				"message": "haven't Password",
			})
			return
		}
		if form.AccNumber == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  1,
				"message": "haven't AccNumber",
			})
			return
		}
	}
	projects := &database.Projects{
		UserName:   form.UserName,
		Password:   form.Password,
		AccNumber:  form.AccNumber,
		ColaAPI:    ColaAPI1,
		StatusJSON: StatusJSON,
	}
	projects.UpdateProjects(form.ID)
	ID, _ := strconv.ParseInt(form.ID, 10, 64)
	data, _ := database.ProjectsCheckID(ID)

	projectsIDInt := strconv.Itoa(int(data.ID))
	projectsIDStr := string(projectsIDInt)
	UsersIDInt := strconv.Itoa(int(data.UsersID))
	UsersIDStr := string(UsersIDInt)
	cache := &CacheValue{
		UsersID:    UsersIDStr,
		ProjectsID: projectsIDStr,
		ColaAPI:    ColaAPI1,
	}
	CacheValues, _ := json.Marshal(&cache)

	Redis.Set(data.Key, string(CacheValues), 0)

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "修改成功",
		"data":    data,
	})
}
