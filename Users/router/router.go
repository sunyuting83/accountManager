package router

import (
	"colaAPI/Users/controller"
	Accounts "colaAPI/Users/controller/Accounts"
	Orders "colaAPI/Users/controller/Orders"
	Users "colaAPI/Users/controller/Users"
	utilsUser "colaAPI/Users/utils"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY, CurrentPath string, FormMemory int64) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = FormMemory << 20
	router.Use(utilsUser.CORSMiddleware(), utilsUser.ThrottleMiddleware())
	userApiV1 := router.Group("/api/v1")
	userApiV1HasKey := router.Group("/api/v1/:key")
	userApiV1.Use(utilsUser.SetConfigMiddleWare(SECRET_KEY, CurrentPath, SECRET_KEY))
	userApiV1HasKey.Use(utilsUser.SetConfigMiddleWare(SECRET_KEY, CurrentPath, SECRET_KEY), utilsUser.UserProjectsMiddleware())
	{
		router.GET("/", controller.Index)
		userApiV1.POST("/Regedit", Users.Regedit)
		userApiV1.POST("/Captcha", Users.Captcha)
		userApiV1.POST("/Login", Users.Sgin)
		userApiV1.PUT("/RePassword", utilsUser.UserVerifyMiddleware(), Users.ResetPassword)
		userApiV1.GET("/CheckLogin", utilsUser.UserVerifyMiddleware(), Users.CheckLogin)
		userApiV1.GET("/GetProducts", utilsUser.UserVerifyMiddleware(), Accounts.GetAccountsList)
		userApiV1.GET("/SearchProducts", utilsUser.UserVerifyMiddleware(), Accounts.SearchProducts)
		userApiV1.POST("/PostOrders", utilsUser.UserVerifyMiddleware(), Orders.PostOrders)
		userApiV1.GET("/GetOrdersList", utilsUser.UserVerifyMiddleware(), Orders.GetOrdersList)
		userApiV1.GET("/GetOrdersDetail", utilsUser.UserVerifyMiddleware(), Orders.GetOrdersDetail)
		userApiV1.POST("/OrderRefund", utilsUser.UserVerifyMiddleware(), Orders.OrderRefund)
	}
	return router
}
