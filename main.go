package main

import (
	Redis "colaAPI/Redis"
	BadgerDB "colaAPI/badger"
	orm "colaAPI/database"
	"colaAPI/router"
	"colaAPI/utils"
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

	confYaml, err := utils.CheckConfig(OS, CurrentPath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	pwd := utils.MD5(strings.Join([]string{confYaml.AdminPWD, confYaml.SECRET_KEY}, ""))
	orm.InitDB(pwd, confYaml)
	Redis.InitRedis(confYaml.Redis.Host, confYaml.Redis.Password, confYaml.Redis.DB)
	gin.SetMode(gin.ReleaseMode)
	// gin.SetMode(gin.DebugMode)
	defer orm.Eloquent.Close()
	defer BadgerDB.BadgerDB.Close()
	defer Redis.MyRedis.Close()
	app := router.InitRouter(confYaml.SECRET_KEY, CurrentPath, confYaml.FormMemory, confYaml.UsersApi.SECRET_KEY)

	// app.Run(strings.Join([]string{":", confYaml.Port}, ""))
	srv := &http.Server{
		Addr:    strings.Join([]string{":", confYaml.Port}, ""),
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
