package controller

import (
	"colaAPI/database"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DeleteList struct {
	Status string `form:"status" json:"status" xml:"status" binding:"required"`
}

func DeleteAccount(c *gin.Context) {
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
	ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)
	Projects, err := database.ProjectsCheckID(ProjectsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	var statusJson []*StatusJSON
	json.Unmarshal([]byte(Projects.StatusJSON), &statusJson)

	var hasPower bool = false

	for _, item := range statusJson {
		if item.Status == form.Status {
			if item.Delete {
				hasPower = true
				break
			}
		}
	}

	if !hasPower {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "状态不支持删除",
		})
		return
	}

	var account *database.Accounts
	account.DeleteAll(projectsID, form.Status)
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "删除成功",
	})
}
