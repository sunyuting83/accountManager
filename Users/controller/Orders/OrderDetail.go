package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOrdersDetail(c *gin.Context) {
	var orderID string = c.DefaultQuery("order_id", "0")
	if orderID == "0" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "参数错误",
		})
		return
	}
	UsersID := utils.GetCurrentUserID(c)
	data, err := database.GetOrdersDetail(orderID, UsersID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "获取数据错误",
		})
		return
	}
	newData := MakeOrderData(data)
	Data := gin.H{
		"status": 0,
		"data":   newData,
	}
	c.JSON(http.StatusOK, Data)
}

type OrderResponseData struct {
	ID          uint
	OrderCode   string
	Status      int
	Coin        float64
	CoinUsersID uint
	Accounts    []*AccountsData
	CreatedAt   int64
	UpdatedAt   int64
}

type AccountsData struct {
	ID        uint
	Status    int
	GameName  string
	Account   string
	Password  string
	Cover     string
	Gold      string
	Multiple  int64
	Diamond   int
	Crazy     int
	Precise   int
	Cold      int
	Price     float64
	Remarks   string
	UpdatedAt int64
}

func MakeOrderData(data *database.Order) *OrderResponseData {
	// fmt.Println(len(*dataList))
	DataList := make([]*AccountsData, len(data.Accounts))
	var result *OrderResponseData = &OrderResponseData{
		ID:          data.ID,
		OrderCode:   data.OrderCode,
		Status:      data.NewStatus,
		Coin:        data.Coin,
		CoinUsersID: data.CoinUsersID,
		Accounts:    DataList,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
	for i, item := range data.Accounts {
		item.Price = utils.Decimal(item.Games.BasePrice + ((item.Games.UnitPrice / float64(item.Games.SingleNumber*100000000)) * float64(item.TodayGold)))
		Account := item.UserName
		Password := item.Password
		if item.SellStatus == 120 {
			Account = utils.ReplaceFromThirdChar(item.UserName, 2)
			Password = ""
		}
		ResponsItems := &AccountsData{
			ID:        item.ID,
			Status:    item.SellStatus,
			GameName:  item.Games.GameName,
			Account:   Account,
			Password:  Password,
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
	return result
}
