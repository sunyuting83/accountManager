package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func UpStatusProjects(c *gin.Context) {
	var form ProjectsID
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	projects, err := database.ProjectsCheckID(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	var (
		NewStatus int    = 1
		FuckStr   string = "锁定"
	)
	if projects.NewStatus == 1 {
		NewStatus = 0
		FuckStr = "解锁"
		projectsIDInt := strconv.Itoa(int(projects.ID))
		UsersID := strconv.Itoa(int(projects.UsersID))
		cache := &CacheValue{
			UsersID:    UsersID,
			ProjectsID: projectsIDInt,
		}
		CacheValues, _ := json.Marshal(&cache)

		Redis.Set(projects.Key, string(CacheValues), 0)
	} else {
		Redis.Delete(projects.Key)
	}
	projects.UpStatusProjects(NewStatus)

	c.JSON(http.StatusOK, gin.H{
		"status":   0,
		"message":  strings.Join([]string{"成功", FuckStr, "后台"}, ""),
		"projects": projects,
	})
}
