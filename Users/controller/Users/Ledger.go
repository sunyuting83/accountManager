package controller

import (
	"colaAPI/Users/database"
	"colaAPI/Users/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLedger(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "50")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
	result := utils.GetTokenUserData(c)
	fmt.Println(result.UserID)
	var bill *database.Bill
	count, err := bill.GetLedgerCount(result.UserID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project",
		})
		return
	}

	ledger, err := database.GetLedger(pageInt, LimitInt, result.UserID)

	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	Data := gin.H{
		"status":  0,
		"message": "获取成功",
		"ledger":  ledger,
		"total":   count,
	}
	c.JSON(http.StatusOK, Data)
}
