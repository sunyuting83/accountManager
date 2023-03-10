package controller

import (
	"colaAPI/database"
	"colaAPI/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ResetPassword(c *gin.Context) {
	var form Login
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

	result := utils.GetTokenUserData(c)

	var admin database.Manager
	user, err := database.CheckUserName(form.UserName)
	if err != nil && err.Error() != "record not found" {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if user.UserName == form.UserName {
		secret_key, _ := c.Get("secret_key")
		SECRET_KEY := secret_key.(string)
		PASSWD := utils.MD5(strings.Join([]string{form.Password, SECRET_KEY}, ""))
		admin.Password = PASSWD
		var (
			data database.Manager
			err  error
		)
		if user.ID == 1 {
			if user.ID == result.UserID {
				data, err = admin.ResetPassword(form.UserName)
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  1,
					"message": "haven't power",
				})
				return
			}
		} else {
			data, err = admin.ResetPassword(form.UserName)
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
