package main

import (
	Redis "colaAPI/Redis"
	BadgerDB "colaAPI/Users/badger"
	orm "colaAPI/Users/database"
	"colaAPI/Users/router"
	"colaAPI/Users/utils"
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
	ConfigFilePath := orm.MakeSqlitePath(CurrentPath)

	confYaml, err := utils.CheckConfig(OS, ConfigFilePath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	orm.InitDB(confYaml)
	Redis.InitRedis(confYaml.Redis.Host, confYaml.Redis.Password, confYaml.Redis.DB)

	// gin.SetMode(gin.DebugMode)
	gin.SetMode(gin.ReleaseMode)
	defer orm.Eloquent.Close()
	defer BadgerDB.BadgerDB.Close()
	defer Redis.MyRedis.Close()
	app := router.InitRouter(confYaml.Users.SECRET_KEY, CurrentPath, confYaml.FormMemory)

	app.Run(strings.Join([]string{":", confYaml.Users.Port}, ""))
}
