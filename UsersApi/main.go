package main

import (
	Redis "colaAPI/Redis"
	BadgerDB "colaAPI/UsersApi/badger"
	"colaAPI/UsersApi/router"
	orm "colaAPI/database"
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
	// gin.SetMode(gin.ReleaseMode)
	Redis.InitRedis(confYaml.Redis.Host, confYaml.Redis.Password, confYaml.Redis.DB)

	gin.SetMode(gin.DebugMode)
	defer orm.Eloquent.Close()
	defer BadgerDB.BadgerDB.Close()
	defer Redis.MyRedis.Close()
	app := router.InitRouter(confYaml.UsersApi.SECRET_KEY, CurrentPath)

	app.Run(strings.Join([]string{":", confYaml.UsersApi.Port}, ""))
}
