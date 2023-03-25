package controller

import (
	"colaAPI/database"
	"colaAPI/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CacheToken struct {
	UserID uint
	Token  string
}

func ProjectsList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	userData := utils.GetTokenUserData(c)

	ManagerData, err := database.GetHasUsersID(userData.UserID)
	// fmt.Println(ManagerData)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "失败",
		})
		return
	}
	HasUsersList := ManagerData.Users
	var UsersIDs []int

	for _, item := range HasUsersList {
		UsersIDs = append(UsersIDs, int(item.ID))
	}

	var (
		projects *database.Projects
		count    int64
		dataList *[]database.Projects
	)
	if userData.UserID == 1 {
		count, err = projects.GetCount()
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "失败",
			})
			return
		}

		dataList, err = database.GetProjectsList(pageInt, LimitInt)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "失败",
			})
			return
		}
	} else {
		count, err = projects.GetCountWithIn(UsersIDs)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "失败",
			})
			return
		}

		dataList, err = database.GetProjectsListWithIn(pageInt, LimitInt, UsersIDs)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "失败",
			})
			return
		}
	}

	Data := gin.H{
		"status": 0,
		"data":   dataList,
		"total":  count,
	}
	c.JSON(http.StatusOK, Data)
}
