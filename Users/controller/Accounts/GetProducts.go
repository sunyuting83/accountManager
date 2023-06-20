package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAccountsList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "50")
	var gameid string = c.DefaultQuery("gameid", "0")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	GameID64, _ := strconv.Atoi(gameid)
	GameID := uint(GameID64)

	var account *database.Accounts
	count, err := account.GetCountWithSellStatus(GameID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project",
		})
		return
	}
	dataList, err := database.GetAccountList(pageInt, LimitInt, GameID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project list",
		})
		return
	}
	DataList := MakeDataList(dataList)
	Data := gin.H{
		"status": 0,
		"data":   DataList,
		"total":  count,
	}
	c.JSON(http.StatusOK, Data)
}

type ResponseDatas struct {
	ID        uint `gorm:"primaryKey"`
	GameID    uint
	GameName  string
	Account   string
	Cover     string
	Gold      string
	Multiple  int64
	Diamond   int
	Crazy     int
	Precise   int
	Cold      int
	Price     float64
	Remarks   string
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func MakeDataList(dataList *[]database.Accounts) []*ResponseDatas {
	// fmt.Println(len(*dataList))
	DataList := make([]*ResponseDatas, len(*dataList))
	for i, item := range *dataList {
		if item.Price <= 0 {
			item.Price = utils.Decimal(item.Games.BasePrice + ((item.Games.UnitPrice / float64(item.Games.SingleNumber*100000000)) * float64(item.TodayGold)))
		}
		ResponsItems := &ResponseDatas{
			ID:        item.ID,
			GameID:    item.GameID,
			GameName:  item.Games.GameName,
			Account:   utils.ReplaceFromThirdChar(item.UserName, 2),
			Cover:     item.Cover,
			Gold:      utils.ConvertNumber(item.TodayGold),
			Multiple:  item.Multiple,
			Diamond:   item.Diamond,
			Crazy:     item.Crazy,
			Precise:   item.Precise,
			Cold:      item.Cold,
			Price:     item.Price,
			Remarks:   item.Remarks,
			UpdatedAt: item.UpdatedAt,
		}
		DataList[i] = ResponsItems
	}
	return DataList
}
