package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type List struct {
	List []int `form:"list" json:"list" xml:"list"  binding:"required"`
}

func PullAccountDrawList(c *gin.Context) {
	var form List
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	tempList := RemoveRepeatedList(form.List)
	if len(tempList) != 0 {
		var acc *database.Accounts

		acc.PullDataUseIn(form.List)
		projectsID, ColaAPI := GetProjects(c)
		if ColaAPI {
			ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)
			Projects, err := database.ProjectsCheckID(ProjectsID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  1,
					"message": err.Error(),
				})
				return
			}
			Projects = &database.Projects{
				UserName:  Projects.UserName,
				Password:  Projects.Password,
				AccNumber: Projects.AccNumber - len(tempList),
			}
			Projects.UpdateProjects(projectsID)
		}
		Data := gin.H{
			"status":  0,
			"message": "提取成功",
		}
		c.JSON(http.StatusOK, Data)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "参数错误",
	})
}

func RemoveRepeatedList(personList []int) (result []int) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < n; j++ {
			if personList[i] == personList[j] {
				repeat = true
				break
			}
		}
		if !repeat && personList[i] != 0 {
			result = append(result, personList[i])
		}
	}
	return
}

func GetProjects(c *gin.Context) (projectsID string, ColaAPI bool) {
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
		ColaAPI = result.ColaAPI
	}
	return
}
