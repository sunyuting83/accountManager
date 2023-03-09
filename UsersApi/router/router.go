package router

import (
	"colaAPI/UsersApi/controller"
	Projects "colaAPI/UsersApi/controller/Projects"
	Users "colaAPI/UsersApi/controller/Users"
	utilsUser "colaAPI/UsersApi/utils"
	"colaAPI/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY, CurrentPath string) *gin.Engine {
	router := gin.Default()
	router.Use(utils.CORSMiddleware())
	userApiV1 := router.Group("/api/v1")
	userApiV1HasKey := router.Group("/api/v1/:key")
	userApiV1.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath))
	userApiV1HasKey.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath))
	{
		router.GET("/", utilsUser.UserVerifyMiddleware(), controller.Index)
		userApiV1.POST("/Login", Users.Sgin)
		userApiV1.PUT("/RePassword", utilsUser.UserVerifyMiddleware(), Users.ResetPassword)
		userApiV1.GET("/CheckLogin", utilsUser.UserVerifyMiddleware(), Users.CheckLogin)
		userApiV1.GET("/ProjectsList", utilsUser.UserVerifyMiddleware(), Projects.ProjectsList)
		userApiV1HasKey.PUT("/UpdateProjects", utilsUser.UserVerifyMiddleware(), Projects.ModifyProjects)
	}

	return router
}
