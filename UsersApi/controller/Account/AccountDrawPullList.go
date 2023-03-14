package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type List struct {
	List []int `form:"list" json:"list" xml:"list"  binding:"required"`
}

func PullAccountDrawList(c *gin.Context) {
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
		var acc *database.Accounts

		acc.PullDataUseIn(form.List)

		Data := gin.H{
			"status":  0,
			"message": "提取成功",
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
