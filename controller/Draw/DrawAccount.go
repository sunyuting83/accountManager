package controller

import (
	controller "colaAPI/controller/Account"
	"colaAPI/database"
	"colaAPI/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DrawIndex(c *gin.Context) {
	userData := utils.GetTokenUserData(c)
	if userData.UserID == 1 {
		gamesList, err := database.GetAllGames()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		for _, item := range gamesList {
			if len(item.Projects) != 0 {
				for _, value := range item.Projects {
					hasStatus := GetPullStatus(value.StatusJSON)
					if len(hasStatus) != 0 {
						var account *database.Accounts
						IDStr := strconv.Itoa(int(value.ID))
						Count, _ := account.GetInCount(IDStr, hasStatus)
						item.Count = item.Count + Count
					}
				}
			}
			item.Projects = make([]database.Projects, 0)
		}

		c.JSON(200, gin.H{
			"status":   0,
			"gamslist": gamesList,
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}

func GetPullStatus(StatusJSON string) []string {
	var statusJson []*controller.StatusJSON
	json.Unmarshal([]byte(StatusJSON), &statusJson)

	var (
		hasStatus []string
	)
	for _, item := range statusJson {
		if item.Pull {
			hasStatus = append(hasStatus, item.Status)
		}
	}
	return hasStatus
}
