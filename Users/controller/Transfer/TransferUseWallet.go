package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransferForUseWallet struct {
	Wallet string `form:"wallet" json:"wallet" xml:"wallet"  binding:"required"`
	Coin   string `form:"coin" json:"coin" xml:"coin"  binding:"required"`
}

func TransferUseWallet(c *gin.Context) {
	var form TransferForUseWallet
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	floatCoin, err := strconv.ParseFloat(form.Coin, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "金额不能 <= 0",
		})
		return
	}
	if floatCoin <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "金额不能 <= 0",
		})
		return
	}
	result := utils.GetTokenUserData(c)

	user, err := database.UserCheckID(int64(result.UserID))
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if user.Coin < floatCoin {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "余额不足",
		})
		return
	}
	// formID, _ := strconv.ParseInt(form.ID, 10, 64)
	formUser, err := database.CheckUserKey(form.Wallet)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if user.ID == formUser.ID {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "不能转账给自己",
		})
		return
	}
	database.UpCoinToCoinUser(user.ID, floatCoin, "-")
	database.UpCoinToCoinUser(formUser.ID, floatCoin, "+")

	d := time.Now()
	dateMonths := d.Format("2006-01-02")

	var bill database.Bill
	bill.Coin = floatCoin
	bill.CoinUsersID = &user.ID
	bill.FormCoinUsersID = &formUser.ID
	bill.NewStatus = 1
	bill.Months = dateMonths
	bill.Insert()

	var billto database.Bill
	billto.Coin = floatCoin
	billto.CoinUsersID = &formUser.ID
	billto.FormCoinUsersID = &user.ID
	billto.NewStatus = 2
	billto.Months = dateMonths
	billto.Insert()

	Data := gin.H{
		"status":    0,
		"message":   "转账成功",
		"form_user": formUser,
		"coin":      form.Coin,
	}
	c.JSON(http.StatusOK, Data)
}
