package controller

import (
	BadgerDB "colaAPI/UsersApi/badger"
	"colaAPI/UsersApi/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetColaToken(c *gin.Context) {
	projectsID, ColaAPI := GetProjectsID(c)
	if ColaAPI {
		token, err := BadgerDB.Get([]byte(projectsID + ".token"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		if len(token) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "get token failed",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "get token succeeded",
			"token":   token,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}

func GetColaAccount(c *gin.Context) {
	projectsID, ColaAPI := GetProjectsID(c)
	if ColaAPI {
		account, err := BadgerDB.Get([]byte(projectsID + ".account"))
		// fmt.Println(account)
		if err != nil && err.Error() == "Key not found" {
			// fmt.Println(err)
			BadgerDB.Set([]byte(projectsID+".account"), []byte("true"))
		}
		if account == "true" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "login status is ing",
			})
			return
		}
		Projects, err := database.ProjectsCheckID(projectsID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "haven't project ID",
			})
			return
		}
		token, _ := BadgerDB.Get([]byte(projectsID + ".token"))
		BadgerDB.Set([]byte(projectsID+".account"), []byte("true"))
		c.JSON(http.StatusOK, gin.H{
			"status":   0,
			"message":  "get token succeeded",
			"username": Projects.UserName,
			"password": Projects.Password,
			"token":    token,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}

func DeleteCola(c *gin.Context) {
	projectsID, ColaAPI := GetProjectsID(c)
	if ColaAPI {
		BadgerDB.Delete([]byte(projectsID + ".account"))
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "delete token succeeded",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}

func SetColaToken(c *gin.Context) {
	projectsID, ColaAPI := GetProjectsID(c)
	if ColaAPI {
		BadgerDB.Delete([]byte(projectsID + ".token"))
		BadgerDB.Set([]byte(projectsID+".account"), []byte("false"))
		token := c.PostForm("token")
		if len(token) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "params token must string",
			})
			return
		}
		BadgerDB.SetWithTTL([]byte(projectsID+".token"), []byte(token), 60*60*24)
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "set token successfully",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}
