package controller

import (
	"colaAPI/UsersApi/database"
	"colaAPI/UsersApi/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PullAccountDrawSelect(c *gin.Context) {
	var form utils.Filter
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	projectsID, ColaAPI := GetProjects(c)
	ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)

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

	rows, _ := database.GetDataUseScopes(form, hasStatus, projectsID)

	if len(rows) != 0 {
		IDs := make([]int, len(rows))
		for _, item := range rows {
			IDs = append(IDs, int(item.ID))
		}
		AdminID := utils.GetCurrentUserID(c)
		user, err := database.UserCheckID(int64(AdminID))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "get admin failed",
			})
			return
		}
		upData, err := database.PullDataUseIn(IDs)
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
			Projects = &database.Projects{
				UserName:  Projects.UserName,
				Password:  Projects.Password,
				AccNumber: Projects.AccNumber - int(len(rows)),
			}
			Projects.UpdateProjects(projectsID)
		}

		Data := gin.H{
			"status":  0,
			"message": "提取成功",
			"data":    rows,
		}
		c.JSON(http.StatusOK, Data)
		return
	}

	empty := make([]string, 0)
	Data := gin.H{
		"status":  1,
		"message": "没有符合条件的帐号",
		"data":    empty,
	}
	c.JSON(http.StatusOK, Data)
}
