package controller

import (
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type StatusJSON struct {
	Status   string `json:"status"`
	Title    string `json:"title"`
	Delete   bool   `json:"delete"`
	CallBack bool   `json:"callback"`
	BackTo   string `json:"backto"`
	Export   bool   `json:"export"`
	Import   bool   `json:"import"`
	Pull     bool   `json:"pull"`
}

func BackToAccount(c *gin.Context) {
	var (
		computid string = c.Query("computid")
		status   string = c.Query("status")
		IsJson   string = c.DefaultQuery("json", "0")
	)
	if len(computid) == 0 {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "haven't params found",
			})
			return
		}
		c.String(200, "参数错误")
		return
	}
	Path := c.Request.URL.Path
	PathList := strings.Split(Path, "/")
	Path = PathList[len(PathList)-1]
	if strings.Contains(Path, "clean") {
		status = GetBackPath(Path)
	}
	if len(status) == 0 {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "haven't params found",
			})
			return
		}
		c.String(200, "参数错误")
		return
	}

	projectsID, _ := GetProjectsID(c)

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
		hasPower     bool = false
		backToStatus string
	)

	for _, item := range statusJson {
		if item.Status == status {
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

	computs, err := database.GetOneComputer(computid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "状态不支持退回",
		})
		return
	}

	backToStatusInt, _ := strconv.Atoi(backToStatus)
	var account *database.Accounts
	account.BackToAcc(projectsID, status, backToStatusInt, computs.ID)
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "退回成功",
	})
}

func GetBackPath(s string) (status string) {
	switch s {
	case "cleanreg":
		status = "1"
	case "cleanpaly":
		status = "3"
	}
	return
}
