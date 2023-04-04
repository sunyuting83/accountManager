package controller

import (
	"colaAPI/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Games struct {
	GameName string `form:"gamename" json:"gamename" xml:"gamename"  binding:"required"`
}

func AddGame(c *gin.Context) {
	var form Games
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if len(form.GameName) < 4 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't username",
		})
		return
	}

	game, err := database.CheckGamesName(form.GameName)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if len(game.GameName) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "游戏已存在",
		})
		return
	}

	game = &database.Games{
		GameName: form.GameName,
	}
	err = game.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "添加成功",
		"data":    game,
	})
}
