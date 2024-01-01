package controller

import (
	BadgerDB "colaAPI/ImageServer/badger"
	"colaAPI/utils"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SetImage(c *gin.Context) {
	name := c.PostForm("name")
	sk := c.PostForm("sk")
	secret_key, _ := c.Get("secret_key")
	SECRET_KEY := secret_key.(string)
	if sk != SECRET_KEY {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "权限验证失败",
		})
		return
	}
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "上传文件失败",
		})
		return
	}

	b, _ := io.ReadAll(file)

	OS := runtime.GOOS
	var SplitString string = "\\"
	if OS == "linux" {
		SplitString = "/"
	}

	fileName := header.Filename
	if strings.Contains(header.Filename, `\`) {
		nameList := strings.Split(header.Filename, `\`)
		Length := len(nameList) - 1
		fileName = nameList[Length]
	}
	fileNameList := strings.Split(fileName, ".")

	if len(fileNameList[0]) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": "6上传文件失败",
		})
		return
	}

	nowTime := time.Now()
	timeStr := nowTime.Format("20060102")
	// fileName = strings.Join([]string{fileNameList[0], "_", gold, ".", fileNameList[1]}, "")
	Path := strings.Join([]string{"upimg", timeStr}, SplitString)
	// fmt.Println(Path)
	imgPath := strings.Join([]string{Path, fileName}, SplitString)

	check := utils.IsExist("upimg")
	if !check {
		err := os.MkdirAll("upimg", 0766)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  0,
				"message": "2上传文件失败",
			})
			return
		}
	}
	checks := utils.IsExist(Path)
	if !checks {
		err := os.MkdirAll(Path, 0766)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  0,
				"message": "2上传文件失败",
			})
			return
		}
	}
	err = os.WriteFile(imgPath, b, 0644)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  0,
			"message": err.Error(),
		})
		return
	}

	var ttl int64 = 60 * 60 * 24 * 30 // ttl以秒为单位
	BadgerDB.SetWithTTL([]byte(name), []byte(imgPath), ttl)

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "上传文件成功",
	})
}
