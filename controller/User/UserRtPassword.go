package controller

import (
	"colaAPI/database"
	"colaAPI/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func UserResetPassword(c *gin.Context) {
	var form User
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	if len(form.UserName) < 4 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't username",
		})
		return
	}
	if len(form.Password) < 8 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't password",
		})
		return
	}

	user, err := database.UserCheckUserName(form.UserName)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if user.UserName == form.UserName {
		secret_key, _ := c.Get("users_secret_key")
		SECRET_KEY := secret_key.(string)
		PASSWD := utils.MD5(strings.Join([]string{form.Password, SECRET_KEY}, ""))
		user.Password = PASSWD
		data, err := user.UserResetPassword(form.UserName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "修改成功",
			"user":    data.UserName,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "用户不存在",
	})
}
