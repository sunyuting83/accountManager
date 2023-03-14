package controller

import (
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAllDateForAccountDraw(c *gin.Context) {

	projectsID := GetProjectsID(c)

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
	hasStatusStr := strings.Join(hasStatus, ",")

	dateList, err := database.GetDateTimeData(projectsID, hasStatusStr, "1")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	Data := gin.H{
		"status":   0,
		"dateList": dateList,
	}
	c.JSON(http.StatusOK, Data)
}
