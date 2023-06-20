package controller

import (
	Accounts "colaAPI/Users/controller/Accounts"
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderList struct {
	List []int `form:"list" json:"list" xml:"list"  binding:"required"`
}

func PostOrders(c *gin.Context) {
	var form OrderList
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	tempList := RemoveRepeatedList(form.List)
	if len(tempList) != 0 {
		accountList, err := database.GetAccountsWithIn(tempList)
		if err != nil && err.Error() == "record not found" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "已经被人抢先了",
			})
			return
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "获取数据失败",
			})
			return
		}

		NewData := MakeDataList(accountList)

		UsersID := utils.GetCurrentUserID(c)
		user, err := database.UserCheckID(int64(UsersID))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "get users failed",
			})
			return
		}
		if user.Coin < NewData.Total {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "余额不足，请先充值",
			})
			return
		}
		ResponseData := gin.H{
			"status": 0,
			"user":   user,
			"data":   NewData,
			"total":  NewData.Total,
			"credit": utils.Decimal(user.Coin - NewData.Total),
		}
		FaileData := filterArray(tempList, accountList)
		if len(FaileData) != 0 {
			failedata, _ := database.GetFailedAccountsWithIn(FaileData)
			newData := Accounts.MakeDataList(failedata)
			ResponseData["FailedData"] = newData
		}
		c.JSON(http.StatusOK, ResponseData)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "参数错误",
	})
}

func RemoveRepeatedList(personList []int) (result []int) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < n; j++ {
			if personList[i] == personList[j] {
				repeat = true
				break
			}
		}
		if !repeat && personList[i] != 0 {
			result = append(result, personList[i])
		}
	}
	return
}

//	func MakeDataList(dataList *[]database.Accounts) float64 {
//		// fmt.Println(len(*dataList))
//		Total := 0.0
//		for _, item := range *dataList {
//			Price := utils.Decimal(item.Games.BasePrice + ((item.Games.UnitPrice / float64(item.Games.SingleNumber*100000000)) * float64(item.TodayGold)))
//			Total += Price
//		}
//		return Total
//	}
type TotalData struct {
	Total      float64
	UniqueItem []UniqueItem
}
type UniqueItem struct {
	Projects database.Projects
	Total    float64
}

func accumulatedPrice(items *[]database.Accounts) map[uint]float64 {
	accumulatedPrice := make(map[uint]float64)
	for _, item := range *items {
		id := item.Projects.ID
		accumulatedPrice[id] += utils.Decimal(item.Games.BasePrice + ((item.Games.UnitPrice / float64(item.Games.SingleNumber*100000000)) * float64(item.TodayGold)))
	}

	return accumulatedPrice
}

// 遍历所有数据，把项目id相同的数据整合成一个数据'UniqueItem',计算每个帐号的价格,项目id相同的总价
func MakeDataList(items *[]database.Accounts) TotalData {
	// 过滤相同项目id的数据并计算总
	accumulatedPrice := accumulatedPrice(items)

	uniqueItems := make([]UniqueItem, 0)
	processedIDs := make(map[uint]bool)
	Total := 0.0
	for _, item := range *items {
		id := item.Projects.ID
		if !processedIDs[id] {
			uniqueItem := UniqueItem{
				Projects: item.Projects,
				Total:    accumulatedPrice[id],
			}
			uniqueItems = append(uniqueItems, uniqueItem)
			Total += accumulatedPrice[id]
			processedIDs[id] = true
		}
	}
	return TotalData{Total: Total, UniqueItem: uniqueItems}
}

func filterArray(arr1 []int, arr2 *[]database.Accounts) []int {
	filtered := make([]int, 0)

	// 创建一个map，用于快速查找arr2中的id值
	idMap := make(map[uint]bool)
	for _, item := range *arr2 {
		idMap[item.ID] = true
	}

	// 遍历arr1，将不在idMap中的值添加到filtered中
	for _, val := range arr1 {
		if !idMap[uint(val)] {
			filtered = append(filtered, val)
		}
	}

	return filtered
}
