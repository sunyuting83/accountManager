package controller

import (
	"colaAPI/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRemarks struct {
	ID      int64  `form:"id" json:"id" xml:"id"  binding:"required"`
	Remarks string `form:"remarks" json:"remarks" xml:"remarks" binding:"required"`
}

func SetUserRemarks(c *gin.Context) {
	var form UserRemarks
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	user, err := database.UserCheckID(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	user.SetUserRemarks(form.Remarks)
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "设置备注成功",
		"id":      user.ID,
	})
}
