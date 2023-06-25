package controller

import (
	"colaAPI/Manager/database"
	"colaAPI/Manager/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SplitProjects struct {
	Percent   float64 `form:"percent" json:"percent" xml:"percent"  binding:"required"`
	ManagerID string  `form:"manager_id" json:"manager_id" xml:"manager_id"  binding:"required"`
}

func SplitForProjects(c *gin.Context) {
	var form SplitProjects
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if form.Percent <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "分成比例必须是数字",
		})
		return
	}
	ManagerID := form.ManagerID
	if len(form.ManagerID) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "分成管理不能为空",
		})
		return
	}
	if !strings.Contains(form.ManagerID, "|||") {
		ManagerID = "1"
	}
	ProportionFloat64 := utils.Decimal(form.Percent)
	err := database.SetSplitProjects(ProportionFloat64, ManagerID)
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
