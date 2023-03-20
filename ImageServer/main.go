package main

import (
	BadgerDB "colaAPI/ImageServer/badger"
	"colaAPI/ImageServer/controller"
	"colaAPI/ImageServer/utils"
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
	gin.SetMode(gin.DebugMode)
	defer BadgerDB.BadgerDB.Close()
	app := gin.Default()
	app.MaxMultipartMemory = confYaml.FormMemory << 20
	{
		app.POST("/set", controller.SetImage)
		app.GET("/image/:path", controller.GetImage)
	}

	app.Run(strings.Join([]string{":", confYaml.Port}, ""))
}
