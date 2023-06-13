package controller

import (
	"bytes"
	BadgerDB "colaAPI/Users/badger"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/freetype"
)

func Captcha(c *gin.Context) {
	// 生成随机验证码
	captcha := generateCaptcha()
	fmt.Println(captcha)
	current_path, _ := c.Get("current_path")
	CurrentPath := current_path.(string)

	// 生成图片并获取Base64编码
	imageBase64 := generateImage(captcha, CurrentPath)
	var ttl int64 = 60 * 5 // ttl以秒为单位
	imageBase64 = strings.Join([]string{"data:image/png;base64", imageBase64}, ",")

	BadgerDB.SetWithTTL([]byte(captcha), []byte(captcha), ttl)
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "获取成功",
		"image":   imageBase64,
	})
}

// 生成随机验证码
func generateCaptcha() string {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	// 生成4位随机数字和字母的验证码
	captcha := ""
	for i := 0; i < 5; i++ {
		if randomGenerator.Intn(2) == 0 {
			// 生成随机数字
			captcha += strconv.Itoa(randomGenerator.Intn(10))
		} else {
			// 生成随机字母（大写）
			captcha += string('A' + randomGenerator.Intn(26))
		}
	}

	return captcha
}

// 生成验证码图片并返回Base64编码
func generateImage(captcha, CurrentPath string) string {
	// 创建图片对象
	img := image.NewRGBA(image.Rect(0, 0, 160, 42))

	// 设置图片背景色为随机颜色
	bgColor := getRandomColor()
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 设置文字颜色为随机颜色
	textColor := getRandomColor()

	// 在图片上绘制验证码文字和干扰线
	fontSize := 20.0
	fontDPI := 72.0
	fontFile := strings.Join([]string{CurrentPath, "font", "arial.ttf"}, "/") // 字体文件路径
	// fontFile := "./arial.ttf" // 字体文件路径

	// 加载字体文件
	fontBytes, err := embedFontFile(fontFile)
	if err != nil {
		panic(err)
	}

	// 创建字体
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	// 创建绘图上下文
	c := freetype.NewContext()
	c.SetDPI(fontDPI)
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.NewUniform(textColor))

	// 绘制验证码文字
	pt := freetype.Pt(10, 30)
	for _, char := range captcha {
		_, err = c.DrawString(string(char), pt)
		if err != nil {
			panic(err)
		}
		pt.X += c.PointToFixed(fontSize * 1.5)
	}

	// 绘制干扰线
	for i := 0; i < 10; i++ {
		startX := rand.Intn(130)
		startY := rand.Intn(42)
		endX := rand.Intn(130)
		endY := rand.Intn(42)
		drawLine(img, startX, startY, endX, endY, textColor)
	}

	// 生成Base64编码
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		panic(err)
	}

	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return imageBase64
}

// 绘制线条
func drawLine(img *image.RGBA, startX, startY, endX, endY int, color color.RGBA) {
	dx := float64(endX - startX)
	dy := float64(endY - startY)

	var steps int
	if abs(dx) > abs(dy) {
		steps = int(abs(dx))
	} else {
		steps = int(abs(dy))
	}

	xIncrement := dx / float64(steps)
	yIncrement := dy / float64(steps)

	x := float64(startX)
	y := float64(startY)

	for i := 0; i <= steps; i++ {
		img.Set(int(x), int(y), color)
		x += xIncrement
		y += yIncrement
	}
}

// 加载字体文件
func embedFontFile(fontFile string) ([]byte, error) {
	fontBytes, err := os.ReadFile(fontFile)
	if err != nil {
		return nil, err
	}

	return fontBytes, nil
}

// 获取随机颜色
func getRandomColor() color.RGBA {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	r := uint8(randomGenerator.Intn(256))
	g := uint8(randomGenerator.Intn(256))
	b := uint8(randomGenerator.Intn(256))

	return color.RGBA{r, g, b, 255}
}

// 计算绝对值
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
