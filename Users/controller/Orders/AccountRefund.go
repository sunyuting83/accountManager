package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderDetail struct {
	ID          string `form:"id" json:"id" xml:"id"  binding:"required"`
	AccountList []int  `form:"account_list" json:"account_list" xml:"account_list"  binding:"required"`
}

func AccountRefund(c *gin.Context) {
	var form OrderDetail
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if form.ID == "0" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "参数错误",
		})
		return
	}
	UsersID := utils.GetCurrentUserID(c)
	data, err := database.GetOrdersDetail(form.ID, UsersID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "获取数据错误",
		})
		return
	}
	AccountsList := RemoveRepeatedList(form.AccountList)
	checkHas := CompareArrays(data.Accounts, AccountsList)
	if !checkHas {
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "非法请求",
			})
			return
		}
	}
	_, err = database.UpAccountsToRefund(AccountsList)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "获取数据错误",
		})
		return
	}
	_, err = database.UpOrdersToRefunding(form.ID, 2, "")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "获取数据错误",
		})
		return
	}
	Data := gin.H{
		"status":  0,
		"message": "申请退单成功，请等待管理员审核",
	}
	c.JSON(http.StatusOK, Data)
}

func CompareArrays(arr1 []database.Accounts, arr2 []int) bool {
	m := make(map[uint]bool)
	for _, num := range arr1 {
		m[num.ID] = true
	}
	for _, num := range arr2 {
		ID := uint(num)
		if !m[ID] {
			return false
		}
	}

	return true
}
