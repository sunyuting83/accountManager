package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetAccountList(c *gin.Context) {
	var Status string = c.DefaultQuery("status", "0")
	var IsJson string = c.DefaultQuery("json", "0")

	projectsID, _ := GetProjectsID(c)

	Projects, err := database.ProjectsCheckID(projectsID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project ID",
		})
		return
	}
	dataList, err := database.GetAccountListWithStatus(projectsID, Status)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project list",
		})
		return
	}
	if IsJson == "1" {
		Data := gin.H{
			"status":   0,
			"data":     dataList,
			"projects": Projects,
		}
		c.JSON(http.StatusOK, Data)
		return
	}
	Data := MakeListData(dataList)
	c.String(200, Data)
}

func MakeListData(dataList *[]database.Accounts) string {
	data := ""
	for _, item := range *dataList {
		CrazyStr := strconv.Itoa(item.Crazy)
		RemarksStr := strings.Join([]string{"----", CrazyStr, "----", item.Remarks}, "")
		Value := strings.Join([]string{item.UserName, RemarksStr}, "")
		Value = strings.Join([]string{Value, "||||"}, "")
		data += Value
	}
	newStr := strings.TrimSuffix(data, "||||")
	return newStr
}
