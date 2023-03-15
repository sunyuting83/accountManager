package router

import (
	"colaAPI/UsersApi/controller"
	Account "colaAPI/UsersApi/controller/Account"
	Projects "colaAPI/UsersApi/controller/Projects"
	Users "colaAPI/UsersApi/controller/Users"
	utilsUser "colaAPI/UsersApi/utils"

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
		router.GET("/", utilsUser.UserVerifyMiddleware(), controller.Index)
		userApiV1.POST("/Login", Users.Sgin)
		userApiV1.PUT("/RePassword", utilsUser.UserVerifyMiddleware(), Users.ResetPassword)
		userApiV1.GET("/CheckLogin", utilsUser.UserVerifyMiddleware(), Users.CheckLogin)
		userApiV1.GET("/ProjectsList", utilsUser.UserVerifyMiddleware(), Projects.ProjectsList)
		userApiV1HasKey.PUT("/UpdateProjects", utilsUser.UserVerifyMiddleware(), Projects.ModifyProjects)
		userApiV1HasKey.GET("/AccountList", utilsUser.UserVerifyMiddleware(), Account.AccountList)
		userApiV1HasKey.POST("/PostAccount", utilsUser.UserVerifyMiddleware(), Account.PostAccount)
		userApiV1HasKey.DELETE("/DeleteAccount", utilsUser.UserVerifyMiddleware(), Account.DeleteAccount)
		userApiV1HasKey.PUT("/GoBackAccount", utilsUser.UserVerifyMiddleware(), Account.GoBackAccount)
		userApiV1HasKey.PUT("/ExportAccount", utilsUser.UserVerifyMiddleware(), Account.ExportAccount)
		userApiV1HasKey.GET("/AccountDrawList", utilsUser.UserVerifyMiddleware(), Account.AccountDrawList)
		userApiV1HasKey.GET("/GetAllDateForDraw", utilsUser.UserVerifyMiddleware(), Account.GetAllDateForAccountDraw)
		userApiV1HasKey.GET("/AccountDrawDateList", utilsUser.UserVerifyMiddleware(), Account.AccountDrawDateList)
		userApiV1HasKey.PUT("/PullDrawList", utilsUser.UserVerifyMiddleware(), Account.PullAccountDrawList)
		userApiV1HasKey.PUT("/PullDrawSelect", utilsUser.UserVerifyMiddleware(), Account.PullAccountDrawSelect)
		userApiV1HasKey.GET("/GetAllDateForDrawed", utilsUser.UserVerifyMiddleware(), Account.GetAllDateForAccountDrawed)
		userApiV1HasKey.GET("/AccountDrawedDateList", utilsUser.UserVerifyMiddleware(), Account.AccountDrawedDateList)
	}

	return router
}
