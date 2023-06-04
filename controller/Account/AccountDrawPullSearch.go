package controller

import (
	"colaAPI/database"
	"colaAPI/utils"
	"encoding/json"
	"net/http"
	"strconv"

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

	var (
		hasStatus []string
	)
	for _, item := range statusJson {
		if item.Pull {
			hasStatus = append(hasStatus, item.Status)
		}
	}
	Ignore := false
	if form.IgnoreSell {
		Ignore = true
	}

	rows, _ := database.GetDataUseScopesB(form, hasStatus, projectsID, Ignore)
	games, _ := database.GetFirstCalc(Projects.GamesID)

	if len(rows) != 0 {
		Data := gin.H{
			"status":   0,
			"message":  "检索成功",
			"data":     rows,
			"projects": Projects,
			"games":    games,
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
