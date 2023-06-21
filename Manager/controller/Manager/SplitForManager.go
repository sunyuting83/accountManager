package controller

import (
	"colaAPI/Manager/database"
	"colaAPI/Users/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SplitManager struct {
	ManagerID  int64   `form:"manager_id" json:"manager_id" xml:"manager_id"  binding:"required"`
	Proportion float64 `form:"proportion" json:"proportion" xml:"proportion"  binding:"required"`
}

func SplitForManager(c *gin.Context) {
	var form SplitManager
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	_, err := database.CheckanagerID(uint(form.ManagerID))
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "查询数据失败",
		})
		return
	}
	if err != nil && err.Error() == "record not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "作者不存在",
		})
		return
	}

	if form.Proportion <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "分成比例必须是数字",
		})
		return
	}
	ProportionFloat64 := utils.Decimal(form.Proportion)
	err = database.SetSplitManager(uint(form.ManagerID), ProportionFloat64)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "查询数据失败",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "查询数据失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "添加分成比例成功",
	})
}
