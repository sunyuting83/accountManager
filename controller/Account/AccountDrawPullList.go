package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/database"
	"colaAPI/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
		AdminID := utils.GetCurrentAdminID(c)
		user, err := database.CheckUserID(AdminID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "get admin failed",
			})
			return
		}
		projectsID, ColaAPI := GetProjects(c)
		ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)
		upData, err := database.PullDataUseIn(form.List)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err,
			})
			return
		}

		upDataJsonStr, _ := json.Marshal(&upData)

		d := time.Now()
		date := d.Format("2006-01-02_15:04:05")
		draw := &database.DrawLogs{
			ProjectsID: uint(ProjectsID),
			Data:       string(upDataJsonStr),
			LogName:    date,
			DrawUser:   user.UserName,
		}
		draw.AddDrawLogs()
		if ColaAPI {
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
			"data":    upData,
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
