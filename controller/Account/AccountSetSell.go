package controller

import (
	"colaAPI/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccountSetSellList(c *gin.Context) {
	var form List
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	tempList := RemoveRepeatedList(form.List)
	if len(tempList) != 0 {
		upData, err := database.SetSellUseIn(form.List)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err,
			})
			return
		}
		Data := gin.H{
			"status":  0,
			"message": "设置成功",
			"data":    upData,
		}
		c.JSON(http.StatusOK, Data)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "参数错误",
	})
}
