package controller

import (
	"colaAPI/Manager/database"
	"colaAPI/Manager/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SendCoin struct {
	UserID    string `form:"userid" json:"userid" xml:"userid"  binding:"required"`
	CoinCount string `form:"coin_count" json:"coin_count" xml:"coin_count"  binding:"required"`
}

func SendCoinToUser(c *gin.Context) {
	var form SendCoin
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if len(form.UserID) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't username",
		})
		return
	}
	if len(form.CoinCount) <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't password",
		})
		return
	}
	CoinFloat, err := strconv.ParseFloat(form.CoinCount, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "币值必须是数字",
		})
		return
	}

	userid, err := database.CheckCoinUserID(form.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "用户不存在",
		})
		return
	}
	// 发放coin到用户
	userid.Coin += CoinFloat
	userid.UpCoinToUsers()

	result := utils.GetTokenUserData(c)

	// 记录发放记录
	var record database.PointsRecord

	record.CoinManagerID = result.UserID
	record.CoinUsersID = userid.ID
	record.Coin = CoinFloat
	record.Insert()

	var issued database.IssuedNumber
	issued.IssuedNumber += CoinFloat
	issued.UpCoinTotalNumber()

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "发放成功",
		"coin":    CoinFloat,
		"UserID":  userid.ID,
	})
}
