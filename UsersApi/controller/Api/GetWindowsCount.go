package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetWindowsCount(c *gin.Context) {
	var (
		number string = c.Query("number")
		IsJson string = c.DefaultQuery("json", "0")
	)

	if len(number) == 0 {
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
	if number == "0" {
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
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "参数错误",
				"project": Projects,
			})
		}
		c.String(200, "参数错误")
		return
	}

	Count, err := database.GetWindowCount(projectsID, number)
	if err != nil {
		if IsJson == "1" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "failed to get accounts count",
			})
		}
		c.String(200, "获取失败")
		return
	}
	if IsJson == "1" {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "获取成功",
			"data":    Count,
		})
		return
	}
	Countstr := strconv.FormatInt(Count, 10)
	c.String(200, Countstr)
}
