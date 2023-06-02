package router

import (
	"colaAPI/UsersApi/controller"
	Account "colaAPI/UsersApi/controller/Account"
	Api "colaAPI/UsersApi/controller/Api"
	Projects "colaAPI/UsersApi/controller/Projects"
	Users "colaAPI/UsersApi/controller/Users"
	utilsUser "colaAPI/UsersApi/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY, CurrentPath string, FormMemory int64) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = FormMemory << 20
	router.Use(utilsUser.CORSMiddleware())
	router.StaticFS("/css", http.Dir("static/css"))
	router.StaticFS("/js", http.Dir("static/js"))
	router.StaticFile("/favicon.ico", "static/favicon.ico")
	router.LoadHTMLGlob("static/index.html")
	userApiV1 := router.Group("/api/v1")
	userApiV1HasKey := router.Group("/api/v1/:key")
	userApiV1.Use(utilsUser.SetConfigMiddleWare(SECRET_KEY, CurrentPath, SECRET_KEY))
	userApiV1HasKey.Use(utilsUser.SetConfigMiddleWare(SECRET_KEY, CurrentPath, SECRET_KEY), utilsUser.UserProjectsMiddleware())
	{
		router.GET("/", controller.Index)
		router.GET("/index", controller.Index)
		router.GET("/project", controller.Index)
		router.GET("/account/:key", controller.Index)
		router.GET("/accountDraw/:key/:type", controller.Index)
		router.GET("/accountDrawed/:key", controller.Index)
		router.GET("/accountFiled/:key", controller.Index)
		userApiV1.POST("/Login", Users.Sgin)
		userApiV1.PUT("/RePassword", utilsUser.UserVerifyMiddleware(), Users.ResetPassword)
		userApiV1.GET("/CheckLogin", utilsUser.UserVerifyMiddleware(), Users.CheckLogin)
		userApiV1.GET("/ProjectsList", utilsUser.UserVerifyMiddleware(), Projects.ProjectsList)
		userApiV1HasKey.PUT("/UpdateProjects", utilsUser.UserVerifyMiddleware(), Projects.ModifyProjects)
		userApiV1HasKey.GET("/AccountList", utilsUser.UserVerifyMiddleware(), Account.AccountList)
		userApiV1HasKey.POST("/PostAccount", utilsUser.UserVerifyMiddleware(), Account.PostAccount)
		userApiV1HasKey.DELETE("/DeleteAccount", utilsUser.UserVerifyMiddleware(), Account.DeleteAccount)
		userApiV1HasKey.PUT("/GoBackAccount", utilsUser.UserVerifyMiddleware(), Account.GoBackAccount)
		userApiV1HasKey.GET("/ExportAccount", utilsUser.UserVerifyMiddleware(), Account.ExportAccount)
		userApiV1HasKey.GET("/AccountDrawList", utilsUser.UserVerifyMiddleware(), Account.AccountDrawList)
		userApiV1HasKey.GET("/GetAllDateForDraw", utilsUser.UserVerifyMiddleware(), Account.GetAllDateForAccountDraw)
		userApiV1HasKey.GET("/AccountDrawDateList", utilsUser.UserVerifyMiddleware(), Account.AccountDrawDateList)
		userApiV1HasKey.PUT("/PullDrawList", utilsUser.UserVerifyMiddleware(), Account.PullAccountDrawList)
		userApiV1HasKey.PUT("/PullDrawSelect", utilsUser.UserVerifyMiddleware(), Account.PullAccountDrawSelect)
		userApiV1HasKey.PUT("/SearchAccountDraw", utilsUser.UserVerifyMiddleware(), Account.SearchAccountDraw)
		userApiV1HasKey.PUT("/SetSellList", utilsUser.UserProjectsMiddleware(), Account.AccountSetSellList)
		userApiV1HasKey.GET("/GetAllDateForDrawed", utilsUser.UserVerifyMiddleware(), Account.GetAllDateForAccountDrawed)
		userApiV1HasKey.GET("/AccountDrawedDateList", utilsUser.UserVerifyMiddleware(), Account.AccountDrawedDateList)
		userApiV1HasKey.GET("/ExportDrawed", utilsUser.UserVerifyMiddleware(), Account.ExportAccountDrawed)
		userApiV1HasKey.GET("/GetFiledList", utilsUser.UserVerifyMiddleware(), Account.GetFiledList)
		userApiV1HasKey.GET("/GetOneFiled", utilsUser.UserVerifyMiddleware(), Account.GetOneFiled)
		userApiV1HasKey.GET("/GetOneAccount", Api.GetOneAccount)
		userApiV1HasKey.GET("/AddAccount", Api.AddAccount)
		userApiV1HasKey.GET("/BackToAccount", Api.BackToAccount)
		userApiV1HasKey.GET("/GoBackAccount", Api.GoBackAccount)
		userApiV1HasKey.GET("/DeleteAccount", Api.DeleteAccount)
		userApiV1HasKey.POST("/PostSetAccount", Api.PostSetAccount)
		userApiV1HasKey.GET("/CheckAccount", Api.CheckAccount)
		userApiV1HasKey.GET("/CheckAccounts", Api.CheckAccountForStatus)
		userApiV1HasKey.GET("/findregone", Api.GetOneAccount)
		userApiV1HasKey.GET("/getplayone", Api.GetOneAccount)
		userApiV1HasKey.GET("/SetAccount", Api.SetAccount)
		userApiV1HasKey.GET("/regonefinished", Api.SetAccount)
		userApiV1HasKey.GET("/playonefinished", Api.SetAccount)
		userApiV1HasKey.GET("/billionfinished", Api.SetAccount)
		userApiV1HasKey.GET("/banfinished", Api.SetAccount)
		userApiV1HasKey.GET("/cleanreg", Api.BackToAccount)
		userApiV1HasKey.GET("/cleanpaly", Api.BackToAccount)
		userApiV1HasKey.GET("/cleanplaycache", Api.GoBackAccount)
		userApiV1HasKey.GET("/GetColaToken", Api.GetColaToken)
		userApiV1HasKey.GET("/ColaLogin", Api.ColaLogin)
		userApiV1HasKey.GET("/ColaAccount", Api.GetColaAccount)
		userApiV1HasKey.PUT("/SetColaToken", Api.SetColaToken)
		userApiV1HasKey.DELETE("/DeleteCola", Api.DeleteCola)
	}
	return router
}
