package router

import (
	"colaAPI/Manager/controller"
	Admin "colaAPI/Manager/controller/Admin"
	Coin "colaAPI/Manager/controller/Coin"
	utils "colaAPI/Manager/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY, CurrentPath string, FormMemory int64) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = FormMemory << 20
	router.Use(utils.CORSMiddleware())
	// router.StaticFS("/css", http.Dir("static/css"))
	// router.StaticFS("/js", http.Dir("static/js"))
	// router.StaticFile("/favicon.ico", "static/favicon.ico")
	// router.LoadHTMLGlob("static/index.html")
	adminapiv1 := router.Group("/admin/api/v1")
	adminApiV1HasKey := router.Group("/admin/api/v1/:key")
	adminapiv1.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath, SECRET_KEY))
	adminApiV1HasKey.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath, SECRET_KEY))
	{
		router.GET("/", controller.Index)
		adminapiv1.POST("/AddAdmin", utils.UserVerifyMiddleware(), Admin.AddAdmin)
		adminapiv1.PUT("/RePassword", utils.UserVerifyMiddleware(), Admin.ResetPassword)
		adminapiv1.DELETE("/DelAdmin", utils.UserVerifyMiddleware(), Admin.DeleteAdmin)
		adminapiv1.GET("/CheckLogin", utils.UserVerifyMiddleware(), Admin.CheckLogin)
		adminapiv1.GET("/AdminList", utils.UserVerifyMiddleware(), Admin.AdminList)
		adminapiv1.PUT("/UpStatus", utils.UserVerifyMiddleware(), Admin.UpStatusAdmin)
		adminapiv1.POST("/AdminLogin", Admin.Sgin)
		adminapiv1.POST("/SendCoinToUser", utils.UserVerifyMiddleware(), Coin.SendCoinToUser)
	}
	return router
}
