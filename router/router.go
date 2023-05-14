package router

import (
	"colaAPI/controller"
	Account "colaAPI/controller/Account"
	Admin "colaAPI/controller/Admin"
	Draw "colaAPI/controller/Draw"
	DrawLog "colaAPI/controller/DrawLog"
	Games "colaAPI/controller/Games"
	Projects "colaAPI/controller/Projects"
	User "colaAPI/controller/User"
	utils "colaAPI/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter make router
func InitRouter(SECRET_KEY, CurrentPath string, FormMemory int64, Users_SECRET_KEY string) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = FormMemory << 20
	router.Use(utils.CORSMiddleware())
	router.StaticFS("/css", http.Dir("static/css"))
	router.StaticFS("/js", http.Dir("static/js"))
	router.StaticFile("/favicon.ico", "static/favicon.ico")
	router.LoadHTMLGlob("static/index.html")
	adminapiv1 := router.Group("/admin/api/v1")
	adminApiV1HasKey := router.Group("/admin/api/v1/:key")
	adminapiv1.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath, Users_SECRET_KEY))
	adminApiV1HasKey.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath, Users_SECRET_KEY))
	{
		router.GET("/", controller.Index)
		router.GET("/adminlist", controller.Index)
		router.GET("/userlist", controller.Index)
		router.GET("/project", controller.Index)
		router.GET("/gameslist", controller.Index)
		router.GET("/account/:key", controller.Index)
		router.GET("/accountDraw/:key/:type", controller.Index)
		router.GET("/accountDrawed/:key", controller.Index)
		router.GET("/accountFiled/:key", controller.Index)
		router.GET("/drawLog/:key", controller.Index)
		router.GET("/drawData/:id", controller.Index)
		router.GET("/AllDraw", controller.Index)
		router.GET("/userProject/:id", controller.Index)
		adminapiv1.POST("/AddAdmin", utils.AdminVerifyMiddleware(), Admin.AddAdmin)
		adminapiv1.PUT("/RePassword", utils.AdminVerifyMiddleware(), Admin.ResetPassword)
		adminapiv1.DELETE("/DelAdmin", utils.AdminVerifyMiddleware(), Admin.DeleteAdmin)
		adminapiv1.GET("/CheckLogin", utils.AdminVerifyMiddleware(), Admin.CheckLogin)
		adminapiv1.GET("/AdminList", utils.AdminVerifyMiddleware(), Admin.AdminList)
		adminapiv1.PUT("/UpStatus", utils.AdminVerifyMiddleware(), Admin.UpStatusAdmin)
		adminapiv1.POST("/AdminLogin", Admin.Sgin)
		adminapiv1.POST("/AddUser", utils.AdminVerifyMiddleware(), User.AddUser)
		adminapiv1.PUT("/RePasswordUser", utils.AdminVerifyMiddleware(), User.UserResetPassword)
		adminapiv1.DELETE("/DelUser", utils.AdminVerifyMiddleware(), User.DeleteUser)
		adminapiv1.GET("/UserList", utils.AdminVerifyMiddleware(), User.UsersList)
		adminapiv1.GET("/UsersAllList", utils.AdminVerifyMiddleware(), User.UsersAllList)
		adminapiv1.PUT("/UpStatusUser", utils.AdminVerifyMiddleware(), User.UpStatusUser)
		adminapiv1.PUT("/SetUserRemarks", utils.AdminVerifyMiddleware(), User.SetUserRemarks)
		adminapiv1.POST("/AddProjects", utils.AdminVerifyMiddleware(), Projects.AddProjects)
		adminapiv1.DELETE("/DelProjects", utils.AdminVerifyMiddleware(), Projects.DeleteProjects)
		adminapiv1.GET("/ProjectsList", utils.AdminVerifyMiddleware(), Projects.ProjectsList)
		adminapiv1.GET("/UserProjectsList", utils.AdminVerifyMiddleware(), Projects.UserProjectsList)
		adminapiv1.PUT("/UpStatusProjects", utils.AdminVerifyMiddleware(), Projects.UpStatusProjects)
		adminapiv1.PUT("/UpdateProjects", utils.AdminVerifyMiddleware(), Projects.ModifyProjects)
		adminapiv1.POST("/AddGame", utils.AdminVerifyMiddleware(), Games.AddGame)
		adminapiv1.DELETE("/DelGame", utils.AdminVerifyMiddleware(), Games.DeleteGame)
		adminapiv1.GET("/GamesList", utils.AdminVerifyMiddleware(), Games.GamesList)
		adminapiv1.GET("/GamesAllList", utils.AdminVerifyMiddleware(), Games.GamesAllList)
		adminapiv1.GET("/DrawData", utils.AdminVerifyMiddleware(), DrawLog.DrawData)
		adminapiv1.GET("/AllCount", utils.AdminVerifyMiddleware(), Draw.DrawIndex)
		adminapiv1.PUT("/DrawSelect", utils.AdminVerifyMiddleware(), Draw.DrawSelect)
		adminapiv1.PUT("/DrawSelectPull", utils.AdminVerifyMiddleware(), Draw.DrawSelectPull)
		adminApiV1HasKey.GET("/DrawList", utils.AdminVerifyMiddleware(), DrawLog.DrawList)
		adminApiV1HasKey.GET("/AccountList", utils.AdminVerifyMiddleware(), Account.AccountList)
		adminApiV1HasKey.POST("/PostAccount", utils.AdminVerifyMiddleware(), Account.PostAccount)
		adminApiV1HasKey.DELETE("/DeleteAccount", utils.AdminVerifyMiddleware(), Account.DeleteAccount)
		adminApiV1HasKey.PUT("/GoBackAccount", utils.AdminVerifyMiddleware(), Account.GoBackAccount)
		adminApiV1HasKey.GET("/ExportAccount", utils.AdminVerifyMiddleware(), Account.ExportAccount)
		adminApiV1HasKey.GET("/AccountDrawList", utils.AdminVerifyMiddleware(), Account.AccountDrawList)
		adminApiV1HasKey.GET("/GetAllDateForDraw", utils.AdminVerifyMiddleware(), Account.GetAllDateForAccountDraw)
		adminApiV1HasKey.GET("/AccountDrawDateList", utils.AdminVerifyMiddleware(), Account.AccountDrawDateList)
		adminApiV1HasKey.PUT("/PullDrawList", utils.AdminVerifyMiddleware(), Account.PullAccountDrawList)
		adminApiV1HasKey.PUT("/PullDrawSelect", utils.AdminVerifyMiddleware(), Account.PullAccountDrawSelect)
		adminApiV1HasKey.GET("/GetAllDateForDrawed", utils.AdminVerifyMiddleware(), Account.GetAllDateForAccountDrawed)
		adminApiV1HasKey.GET("/AccountDrawedDateList", utils.AdminVerifyMiddleware(), Account.AccountDrawedDateList)
		adminApiV1HasKey.GET("/ExportDrawed", utils.AdminVerifyMiddleware(), Account.ExportAccountDrawed)
		adminApiV1HasKey.GET("/GetFiledList", utils.AdminVerifyMiddleware(), Account.GetFiledList)
		adminApiV1HasKey.GET("/GetOneFiled", utils.AdminVerifyMiddleware(), Account.GetOneFiled)
	}

	return router
}
