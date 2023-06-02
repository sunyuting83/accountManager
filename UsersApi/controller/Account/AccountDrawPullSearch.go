package controller

import (
	"colaAPI/UsersApi/database"
	"colaAPI/UsersApi/utils"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchAccountDraw(c *gin.Context) {
	var form utils.Filter
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	projectsID, _ := GetProjects(c)

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
		Data := gin.H{
			"status":  0,
			"message": "检索成功",
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
