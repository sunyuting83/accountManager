package controller

import (
	"colaAPI/Users/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGamesList(c *gin.Context) {
	dataList, err := database.GetGamesList()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't games list",
		})
		return
	}
	DataList := MakeGamesDataList(dataList)
	Data := gin.H{
		"status": 0,
		"data":   DataList,
	}
	c.JSON(http.StatusOK, Data)
}

type GamesResponseDatas struct {
	ID       uint
	GameName string
}

func MakeGamesDataList(dataList *[]database.Games) []*GamesResponseDatas {
	// fmt.Println(len(*dataList))
	DataList := make([]*GamesResponseDatas, len(*dataList))
	for i, item := range *dataList {
		ResponsItems := &GamesResponseDatas{
			ID:       item.ID,
			GameName: item.GameName,
		}
		DataList[i] = ResponsItems
	}
	return DataList
}
