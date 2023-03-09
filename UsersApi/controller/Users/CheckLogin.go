package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Key string `uri:"key" binding:"required"`
}
type CacheValue struct {
	UsersID    string `json:"UsersID"`
	ProjectsID string `json:"ProjectsID"`
}

func CheckLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "登陆中",
	})
}
