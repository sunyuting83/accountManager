package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"math"
	"net/http"
	"strconv"
	"strings"

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
		ProjectsPercent := 7.00
		var CoinManager []string
		SpiltPercent, err := database.GetSplittedPercent()
		if err == nil {
			ProjectsPercent = SpiltPercent.Percent
			if strings.Contains(SpiltPercent.Manager, "|||") {
				CoinManager = strings.Split(SpiltPercent.Manager, "|||")
			} else {
				CoinManager[0] = SpiltPercent.Manager
			}
		}
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
		if len(*accountList) == 0 {
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

		NewID := make([]string, 0)
		for _, item := range *accountList {
			idStr := strconv.Itoa(int(item.ID))
			NewID = append(NewID, idStr)
		}
		database.UpAccountsWithIn(NewID)

		Accounts := strings.Join(NewID, "|||")

		OrderCode := MakeOrderCode(UsersID)
		var order database.Order
		order.Coin = NewData.Total
		order.OrderCode = OrderCode
		order.CoinUsersID = UsersID
		order.AccountsID = Accounts
		order.Insert()

		// 获取所有作者和工作室ID,并计算分成
		for _, item := range NewData.UniqueItem {
			ManagerPercent, _ := database.CheckSplitManagerID(item.Projects.Users.ManagerID)
			CoinManagerPercent := 10.00 - (ProjectsPercent + ManagerPercent.Percent)
			ProjectsCoin, ManagerCoin, CoinManagerCoin := splitAmount(item.Total, ProjectsPercent, ManagerPercent.Percent, CoinManagerPercent)
			remainder := utils.Decimal(item.Total - (roundUpToTwoDecimalPlaces(ProjectsCoin) + ManagerCoin + CoinManagerCoin))
			// fmt.Println(len(CoinManager))
			ProjectsCoin = utils.Decimal(ProjectsCoin)
			if len(CoinManager) > 1 {
				CoinManagerCoin = CoinManagerCoin / float64(len(CoinManager))
			}
			if remainder > 0 {
				ProjectsCoin = ProjectsCoin + remainder
			}
			database.UpCoinToCoinManager(CoinManagerCoin, CoinManager)
			database.UpCoinToManager(ManagerCoin, item.Projects.Users.ManagerID)
			database.UpCoinToUsers(ProjectsCoin, item.Projects.UsersID)
			// fmt.Println(item.Total, ProjectsCoin, ManagerCoin, CoinManagerCoin, remainder)
			var blockchain database.BlockChain
			// UsersID 工作室
			blockchain.UsersID = item.Projects.UsersID
			blockchain.UsersCoin = ProjectsCoin
			blockchain.UsersPercent = ProjectsPercent
			// ManagerID 作者
			blockchain.ManagerID = item.Projects.Users.ManagerID
			blockchain.ManagerCoin = ManagerCoin
			blockchain.ManagerPercent = ManagerPercent.Percent
			// CoinManagerIDs 分润帐号 CoinManagerCoin 单帐号分润金额 CoinManagerPercent 分润总比例 需要除以 len(CoinManagerIDs)
			blockchain.CoinManagerIDs = SpiltPercent.Manager
			blockchain.CoinManagerCoin = CoinManagerCoin
			blockchain.CoinManagerPercent = CoinManagerPercent
			blockchain.Insert()
		}

		ResponseData := gin.H{
			"status": 0,
			"user":   user,
			"total":  NewData.Total,
			"credit": utils.Decimal(user.Coin - NewData.Total),
		}
		FaileData := filterArray(tempList, accountList)
		if len(FaileData) != 0 {
			failedata, _ := database.GetFailedAccountsWithIn(FaileData)
			newData := MakeDataList(failedata)
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

func splitAmount(amount, Projects, Manager, CoinManager float64) (float64, float64, float64) {
	share1 := (Projects / 10) * amount
	share2 := (Manager / 10) * amount
	share3 := (CoinManager / 10) * amount

	return share1, share2, share3
}

func roundUpToTwoDecimalPlaces(f float64) float64 {
	return math.Floor(f*100) / 100
}

func MakeOrderCode(id uint) string {
	idStr := strconv.Itoa(int(id))
	if id < 10 {
		idStr = "0" + idStr
	}
	before := "2023001"
	datetime := utils.GetDateTimeStr()
	orderCode := strings.Join([]string{before, idStr, datetime}, "")
	return orderCode
}
