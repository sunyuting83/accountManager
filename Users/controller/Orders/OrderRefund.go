package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderID struct {
	ID      string `form:"id" json:"id" xml:"id"  binding:"required"`
	Remarks string `form:"remarks" json:"remarks" xml:"remarks"`
}

func OrderRefund(c *gin.Context) {
	var form OrderID
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
	_, err := database.GetOrdersDetailForRefund(form.ID, UsersID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "获取数据错误",
		})
		return
	}
	_, err = database.UpOrdersToRefunding(form.ID, 1, form.Remarks)
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
