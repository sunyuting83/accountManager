package controller

import (
	"colaAPI/database"
	"colaAPI/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Projects struct {
	ID           uint
	UserName     string
	Remarks      string
	ProjectsName string
	Accounts     []database.Accounts
	Count        int
}

func DrawSelect(c *gin.Context) {
	userData := utils.GetTokenUserData(c)
	if userData.UserID == 1 {
		var form utils.Filter
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		if form.GameID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": "game id must",
			})
			return
		}
		gams, err := database.GetGame(form.GameID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": "game id must",
			})
			return
		}
		dataList := make([]*Projects, 0)

		if len(gams.Projects) != 0 {
			for _, value := range gams.Projects {
				hasStatus := GetPullStatus(value.StatusJSON)
				if len(hasStatus) != 0 {
					projectsID := strconv.FormatInt(int64(value.ID), 10)
					value.Accounts, _ = database.GetDataUseScopes1(form, hasStatus, projectsID)
					Count := len(value.Accounts)
					if Count != 0 {
						dataList = append(dataList, &Projects{
							ID:           value.ID,
							UserName:     value.Users.UserName,
							ProjectsName: value.ProjectsName,
							Remarks:      value.Users.Remarks,
							Accounts:     value.Accounts,
							Count:        Count,
						})
					}
				}
			}
		}
		gams.Projects = make([]database.Projects, 0)

		if len(dataList) != 0 {
			Data := gin.H{
				"status":  0,
				"message": "查询成功",
				"data":    dataList,
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
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}
