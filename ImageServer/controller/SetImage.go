package controller

import (
	BadgerDB "colaAPI/ImageServer/badger"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetImage(c *gin.Context) {
	name := c.PostForm("name")
	sk := c.PostForm("sk")
	secret_key, _ := c.Get("secret_key")
	SECRET_KEY := secret_key.(string)
	if sk != SECRET_KEY {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "权限验证失败",
		})
		return
	}
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "上传文件失败",
		})
		return
	}

	b, _ := io.ReadAll(file)
	var ttl int64 = 60 * 60 * 24 * 30 // ttl以秒为单位
	BadgerDB.SetWithTTL([]byte(name), b, ttl)

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "上传文件成功",
	})
}
