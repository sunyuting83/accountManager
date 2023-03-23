package controller

import (
	colaapi "colaAPI/UsersApi/ColaAPI"
	BadgerDB "colaAPI/UsersApi/badger"
	"colaAPI/UsersApi/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ColaLogin(c *gin.Context) {

	projectsID, ColaAPI := GetProjectsID(c)

	if ColaAPI {
		Projects, err := database.ProjectsCheckID(projectsID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "haven't project ID",
			})
			return
		}
		token, err := BadgerDB.Get([]byte(projectsID + ".token"))
		if err != nil || len(token) == 0 {
			token, err = colaapi.Login(Projects.UserName, Projects.Password)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status":  1,
					"message": "login failed",
				})
				return
			}
			if len(token) == 0 {
				c.JSON(http.StatusOK, gin.H{
					"status":  1,
					"message": "login failed",
				})
				return
			}
			BadgerDB.SetWithTTL([]byte(projectsID+".token"), []byte(token), 60*60*24)
			c.JSON(http.StatusOK, gin.H{
				"status":  0,
				"message": "get token successfully",
				"token":   token,
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": "haven't powered up yet",
	})
}
