package router

import (
	"colaAPI/Users/controller"
	Users "colaAPI/Users/controller/Users"
	utilsUser "colaAPI/Users/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY, CurrentPath string, FormMemory int64) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = FormMemory << 20
	router.Use(utilsUser.CORSMiddleware())
	userApiV1 := router.Group("/api/v1")
	userApiV1HasKey := router.Group("/api/v1/:key")
	userApiV1.Use(utilsUser.SetConfigMiddleWare(SECRET_KEY, CurrentPath, SECRET_KEY))
	userApiV1HasKey.Use(utilsUser.SetConfigMiddleWare(SECRET_KEY, CurrentPath, SECRET_KEY), utilsUser.UserProjectsMiddleware())
	{
		router.GET("/", controller.Index)
		userApiV1.POST("/Regedit", Users.Regedit)
		userApiV1.POST("/Captcha", Users.Captcha)
	}
	return router
}
