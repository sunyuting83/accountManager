package controller

import (
	"colaAPI/database"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type GameCalc struct {
	ID           string `form:"id" json:"id" xml:"id"  binding:"required"`
	BasePrice    string `form:"bprice" json:"bprice" xml:"bprice"  binding:"required"`
	UnitPrice    string `form:"uprice" json:"uprice" xml:"uprice"  binding:"required"`
	SingleNumber int64  `form:"number" json:"number" xml:"number"  binding:"required"`
}

func ModifyCalc(c *gin.Context) {
	var form GameCalc
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if len(form.ID) <= 0 || strings.Contains(form.ID, "-") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't ID",
		})
		return
	}
	BasePrice, err := strconv.ParseFloat(form.BasePrice, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't Base-Price",
		})
		return
	}
	UnitPrice, err := strconv.ParseFloat(form.UnitPrice, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't Unit-Price",
		})
		return
	}
	if UnitPrice <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't Unit-Price",
		})
		return
	}
	if form.SingleNumber <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't Number",
		})
		return
	}
	fmt.Println(form.ID)
	GameIDInt, _ := strconv.ParseInt(form.ID, 10, 64)
	modifyCalc, err := database.GameCheckID(GameIDInt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't Number",
		})
		return
	}
	modifyCalc.BasePrice = BasePrice
	modifyCalc.UnitPrice = UnitPrice
	modifyCalc.SingleNumber = form.SingleNumber
	// modifyCale := &database.Games{
	// 	BasePrice:    BasePricpageInt, _ := strconv.Atoi(page)e,
	// 	UnitPrice:    UnitPrice,
	// 	SingleNumber: form.SingleNumber,
	// }
	modifyCalc.UpdateGames()

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "添加成功",
	})
}
