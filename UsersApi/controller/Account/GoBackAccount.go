package controller

import (
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GoBackAccount(c *gin.Context) {
	var form DeleteList
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	projectsID := GetProjectsID(c)
	ProjectsID, _ := strconv.Atoi(projectsID)

	Projects, err := database.ProjectsCheckID(int64(ProjectsID))
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
		hasPower     bool = false
		backToStatus string
	)

	for _, item := range statusJson {
		if item.Status == form.Status {
			if item.CallBack {
				hasPower = true
				backToStatus = item.BackTo
				break
			}
		}
	}

	if !hasPower {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "状态不支持退回",
		})
		return
	}

	var account *database.Accounts
	account.BackTo(projectsID, form.Status, backToStatus)
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "退回成功",
	})
}
