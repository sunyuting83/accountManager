package controller

import (
	"colaAPI/Manager/database"
	"colaAPI/Manager/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type RtPassword struct {
	Password string `form:"password" json:"password" xml:"password"  binding:"required"`
}

func ResetPassword(c *gin.Context) {
	var form RtPassword
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
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

	result := utils.GetTokenUserData(c)

	var admin database.CoinManager
	user, err := database.CoinManagerCheckID(int64(result.UserID))
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if user.ID == result.UserID {
		secret_key, _ := c.Get("secret_key")
		SECRET_KEY := secret_key.(string)
		PASSWD := utils.MD5(strings.Join([]string{form.Password, SECRET_KEY}, ""))
		admin.Password = PASSWD
		var (
			data database.CoinManager
			err  error
		)
		if user.ID == 1 {
			if user.ID == result.UserID {
				data, err = admin.ResetPassword(user.UserName)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  1,
					"message": "haven't power",
				})
				return
			}
		} else {
			data, err = admin.ResetPassword(user.UserName)
		}
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
		"message": "管理员不存在",
	})
}
