package main

import (
	Redis "colaAPI/Redis"
	BadgerDB "colaAPI/Users/badger"
	orm "colaAPI/Users/database"
	"colaAPI/Users/router"
	"colaAPI/Users/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	OS := runtime.GOOS
	CurrentPath, _ := utils.GetCurrentPath()
	ConfigFilePath := orm.MakeSqlitePath(CurrentPath)
	utils.CheckGeoIP(OS, CurrentPath)

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

	// app.Run(strings.Join([]string{":", confYaml.Users.Port}, ""))
	srv := &http.Server{
		Addr:    strings.Join([]string{":", confYaml.Users.Port}, ""),
		Handler: app,
	}
	fmt.Printf("listen port %s\n", srv.Addr)
	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
