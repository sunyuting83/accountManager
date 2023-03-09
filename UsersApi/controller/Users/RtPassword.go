package controller

import (
	BadgerDB "colaAPI/UsersApi/badger"
	"colaAPI/UsersApi/database"
	"colaAPI/utils"
	"encoding/json"
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
	token := c.GetHeader("Authorization")

	secret_key, _ := c.Get("secret_key")
	SECRET_KEY := secret_key.(string)
	token = token[7:]
	AEStoken, err := utils.DecryptByAes(token, []byte(SECRET_KEY))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't token",
		})
		return
	}
	Token, err := BadgerDB.GetToken(AEStoken)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	var result *CacheToken
	json.Unmarshal(Token, &result)
	user, err := database.UserCheckID(int64(result.UserID))
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
		user.Password = PASSWD
		user.UserResetPassword(int64(result.UserID))
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
			"user":    user.UserName,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "用户不存在",
	})
}
