package controller

import (
	"colaAPI/Users/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartList struct {
	Cart []int `form:"cart" json:"cart" xml:"cart"  binding:"required"`
}

func GetCart(c *gin.Context) {
	var form CartList
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	tempList := RemoveRepeatedList(form.Cart)
	if len(tempList) != 0 {
		cartData, err := database.GetCartWithIn(tempList)
		if err != nil && err.Error() != "record not found" {
			c.JSON(200, gin.H{
				"status":  1,
				"message": "发生错误",
			})
			return
		}
		DataList := MakeDataList(cartData)
		Data := gin.H{
			"status": 0,
			"data":   DataList,
		}
		c.JSON(http.StatusOK, Data)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "参数错误",
	})
}

func RemoveRepeatedList(personList []int) (result []int) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < n; j++ {
			if personList[i] == personList[j] {
				repeat = true
				break
			}
		}
		if !repeat && personList[i] != 0 {
			result = append(result, personList[i])
		}
	}
	return
}
