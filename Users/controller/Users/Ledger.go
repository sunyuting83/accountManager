package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLedger(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "50")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	result := utils.GetTokenUserData(c)
	var bill *database.Bill
	count, err := bill.GetLedgerCount(result.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project",
		})
		return
	}

	ledger, err := database.GetLedger(pageInt, LimitInt, result.UserID)

	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	newData := MakeDataList(ledger)

	Data := gin.H{
		"status":  0,
		"message": "获取成功",
		"ledger":  newData,
		"total":   count,
	}
	c.JSON(http.StatusOK, Data)
}

type ResponseDatas struct {
	ID            uint
	Coin          float64
	Status        int
	CreatedAt     int64
	FormCoinUsers string
	OrderCode     string
	StatusName    string
}

func MakeDataList(ledger *[]database.Bill) []*ResponseDatas {
	// fmt.Println(len(*dataList))
	DataList := make([]*ResponseDatas, len(*ledger))
	for i, item := range *ledger {
		UserName := ""
		OrderNumber := ""
		if item.NewStatus == 1 || item.NewStatus == 2 {
			UserName = item.FormCoinUsers.UserName
		}
		if item.NewStatus == 3 {
			OrderNumber = item.Order.OrderCode
		}
		ResponsItems := &ResponseDatas{
			ID:            item.ID,
			Coin:          item.Coin,
			Status:        item.NewStatus,
			CreatedAt:     item.CreatedAt,
			FormCoinUsers: UserName,
			OrderCode:     OrderNumber,
			StatusName:    MakeStatusName(item.NewStatus),
		}
		DataList[i] = ResponsItems
	}
	return DataList
}

func MakeStatusName(status int) (StatusName string) {
	switch status {
	case 0:
		StatusName = "充值 + "
	case 1:
		StatusName = "转账 - "
	case 2:
		StatusName = "转账 + "
	case 3:
		StatusName = "消费 - "
	default:
		StatusName = ""
	}
	return
}
