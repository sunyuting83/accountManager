package router

import (
	"colaAPI/controller"
	Account "colaAPI/controller/Account"
	Admin "colaAPI/controller/Admin"
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
	adminapiv1.Use(utils.SetConfigMiddleWare(SECRET_KEY, CurrentPath, Users_SECRET_KEY))
	{
		router.GET("/", controller.Index)
		router.GET("/adminlist", controller.Index)
		router.GET("/userlist", controller.Index)
		router.GET("/project", controller.Index)
		router.GET("/account/:id", controller.Index)
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
		adminapiv1.GET("/aaa", utils.AdminVerifyMiddleware(), controller.Index)
		adminapiv1.POST("/AddProjects", utils.AdminVerifyMiddleware(), Projects.AddProjects)
		adminapiv1.DELETE("/DelProjects", utils.AdminVerifyMiddleware(), Projects.DeleteProjects)
		adminapiv1.GET("/ProjectsList", utils.AdminVerifyMiddleware(), Projects.ProjectsList)
		adminapiv1.PUT("/UpStatusProjects", utils.AdminVerifyMiddleware(), Projects.UpStatusProjects)
		adminapiv1.PUT("/UpdateProjects", utils.AdminVerifyMiddleware(), Projects.ModifyProjects)
		adminapiv1.GET("/AccountList", utils.AdminVerifyMiddleware(), Account.AccountList)
	}

	return router
}
