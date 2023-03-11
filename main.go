package main

import (
	Redis "colaAPI/Redis"
	BadgerDB "colaAPI/badger"
	orm "colaAPI/database"
	"colaAPI/router"
	"colaAPI/utils"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	OS := runtime.GOOS
	CurrentPath, _ := utils.GetCurrentPath()

	confYaml, err := utils.CheckConfig(OS, CurrentPath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	pwd := utils.MD5(strings.Join([]string{confYaml.AdminPWD, confYaml.SECRET_KEY}, ""))
	orm.InitDB(pwd, confYaml)
	Redis.InitRedis(confYaml.Redis.Host, confYaml.Redis.Password, confYaml.Redis.DB)
	// gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	defer orm.Eloquent.Close()
	defer BadgerDB.BadgerDB.Close()
	defer Redis.MyRedis.Close()
	app := router.InitRouter(confYaml.SECRET_KEY, CurrentPath, confYaml.FormMemory, confYaml.UsersApi.SECRET_KEY)

	app.Run(strings.Join([]string{":", confYaml.Port}, ""))
}
