package controller

import (
	"colaAPI/UsersApi/database"
	"colaAPI/UsersApi/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchAccountDraw(c *gin.Context) {
	var form utils.SearchFilter
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
	Ignore := false
	if form.IgnoreSell {
		Ignore = true
	}

	rows, _ := database.GetDataUseScopesB(form, projectsID, Ignore)
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
