package router

import (
	"colaAPI/UsersApi/controller"
	Users "colaAPI/UsersApi/controller/Users"
	utilsUser "colaAPI/UsersApi/utils"
	"colaAPI/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY, CurrentPath string) *gin.Engine {
	router := gin.Default()
	router.Use(utils.CORSMiddleware())
	userapiiv1 := router.Group("/api/v1/:key")
	userapiiv1.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath))
	{
		router.GET("/", utilsUser.UserVerifyMiddleware(), controller.Index)
		userapiiv1.POST("/Login", Users.Sgin)
		userapiiv1.GET("/", utils.UserProjectsMiddleware(), controller.Index)
	}

	return router
}
