package controller

import (
	"bytes"
	BadgerDB "colaAPI/Users/badger"
	"colaAPI/Users/utils"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
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
	// fmt.Println(captcha)
	current_path, _ := c.Get("current_path")
	CurrentPath := current_path.(string)

	// 生成图片并获取Base64编码
	image := generateImage(captcha, CurrentPath)
	var ttl int64 = 60 * 5 // ttl以秒为单位
	// imageBase64 = strings.Join([]string{"data:image/png;base64", imageBase64}, ",")
	// fmt.Println(captcha)
	VCode := utils.ConvertToUpperCase(captcha)
	BadgerDB.SetWithTTL([]byte(VCode), []byte(VCode), ttl)
	c.Writer.Write(image)
}

func generateRandomLetter() string {
	rand.Seed(time.Now().UnixNano())

	// 生成随机数字 0 或 1
	random := rand.Intn(2)

	// 生成随机字母（大写或小写）
	if random == 0 {
		return string('A' + rand.Intn(26))
	} else {
		return string('a' + rand.Intn(26))
	}
}

// 生成随机验证码
func generateCaptcha() string {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	captcha := ""
	for i := 0; i < 6; i++ {
		// fmt.Println(randomGenerator.Intn(2))
		if randomGenerator.Intn(2) == 0 {
			captcha += strconv.Itoa(randomGenerator.Intn(10))
		} else {
			captcha += generateRandomLetter()
		}
	}

	return captcha
}

// 生成验证码图片并返回Base64编码
func generateImage(captcha, CurrentPath string) []byte {
	// 创建图片对象
	width, height := 160, 42
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 设置图片背景色为随机颜色
	bgColor := getRandomColor()
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 在图片上绘制验证码文字和干扰线
	// fontSize := randomFontSize()
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
	c.SetClip(img.Bounds())
	c.SetDst(img)

	// 绘制验证码文字
	pt := freetype.Pt(10, 30)
	for _, char := range captcha {
		fontSize := randomFontSize()
		c.SetFontSize(fontSize)
		c.SetSrc(image.NewUniform(getRandomColor()))

		_, err = c.DrawString(string(char), pt)
		if err != nil {
			panic(err)
		}
		pt.X += c.PointToFixed(randomCharacterSpacing(fontSize))
	}

	// 绘制干扰线
	for i := 0; i < len(captcha)*2; i++ {
		startX := rand.Intn(width)
		startY := rand.Intn(height)
		endX := rand.Intn(width)
		endY := rand.Intn(height)
		drawLine(img, startX, startY, endX, endY, getRandomColor())
	}

	// 生成Base64编码
	var buf bytes.Buffer
	png.Encode(&buf, img)

	// imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return buf.Bytes()

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

// 生成随机字符间距
func randomCharacterSpacing(fontsize float64) float64 {
	minSpacing := -2
	maxSpacing := 1
	spacing := rand.Intn(maxSpacing-minSpacing+1) + minSpacing
	return fontsize + float64(spacing)
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

func randomFontSize() float64 {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)
	fontSize := randomGenerator.Intn(15) + 18 // 随机生成18-32之间的整数

	return float64(fontSize)
}
