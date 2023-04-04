package controller

import (
	"colaAPI/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GameID struct {
	ID int64 `form:"id" json:"id" xml:"id"  binding:"required"`
}

func DeleteGame(c *gin.Context) {
	var form GameID
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	game, err := database.GameCheckID(form.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	game.DeleteOne(form.ID)
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "成功删除用户",
		"id":      game.ID,
	})
}
