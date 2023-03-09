package controller

import (
	BadgerDB "colaAPI/UsersApi/badger"
	"colaAPI/UsersApi/database"
	"colaAPI/UsersApi/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Key string `uri:"key" binding:"required"`
}
type CacheValue struct {
	UsersID    string `json:"UsersID"`
	ProjectsID string `json:"ProjectsID"`
}
type CacheToken struct {
	UserID uint
	Token  string
}

func ProjectsList(c *gin.Context) {
	var page string = c.DefaultQuery("page", "0")
	var Limit string = c.DefaultQuery("limit", "100")
	pageInt, _ := strconv.Atoi(page)
	LimitInt, _ := strconv.Atoi(Limit)
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
	var projects *database.Projects
	count, err := projects.GetCount(int64(result.UserID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project Count",
		})
		return
	}
	dataList, err := database.GetProjectsList(int64(result.UserID), pageInt, LimitInt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "haven't project list",
		})
		return
	}
	Data := gin.H{
		"status": 0,
		"data":   dataList,
		"total":  count,
	}
	c.JSON(http.StatusOK, Data)
}
