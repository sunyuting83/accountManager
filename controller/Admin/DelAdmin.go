package controller

import (
	"colaAPI/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID int64 `form:"id" json:"id" xml:"id"  binding:"required"`
}

func DeleteAdmin(c *gin.Context) {
	var form User
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	user, err := database.CheckID(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if user.ID != 1 {
		user.DeleteOne(form.ID)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "成功删除管理员",
		"id":      user.ID,
	})
}
