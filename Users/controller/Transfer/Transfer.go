package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TransferForUser struct {
	ID   string  `form:"id" json:"id" xml:"id"  binding:"required"`
	Coin float64 `form:"coin" json:"coin" xml:"coin"  binding:"required"`
}

func Transfer(c *gin.Context) {
	var form TransferForUser
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
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
	if user.Coin < form.Coin {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "余额不足",
		})
		return
	}
	formID, _ := strconv.ParseInt(form.ID, 10, 64)
	formUser, err := database.UserCheckID(formID)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	database.UpCoinToCoinUser(user.ID, form.Coin, "-")
	database.UpCoinToCoinUser(formUser.ID, form.Coin, "+")

	d := time.Now()
	dateMonths := d.Format("2006-01-02")

	var bill database.Bill
	bill.Coin = form.Coin
	bill.CoinUsersID = &user.ID
	bill.FormCoinUsersID = &formUser.ID
	bill.NewStatus = 1
	bill.Months = dateMonths
	bill.Insert()

	bill.Coin = form.Coin
	bill.CoinUsersID = &formUser.ID
	bill.FormCoinUsersID = &user.ID
	bill.NewStatus = 2
	bill.Months = dateMonths
	bill.Insert()

	Data := gin.H{
		"status":    0,
		"message":   "转账成功",
		"form_user": formUser,
		"coin":      form.Coin,
	}
	c.JSON(http.StatusOK, Data)
}
