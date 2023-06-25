package utils

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port       string      `yaml:"port"`
	SECRET_KEY string      `yaml:"SECRET_KEY"`
	AdminPWD   string      `yaml:"AdminPWD"`
	FormMemory int64       `yaml:"FormMemory"`
	Database   Database    `yaml:"Database"`
	UsersApi   UsersApi    `yaml:"UsersApi"`
	ManagerApi ManagerApi  `yaml:"ManagerApi"`
	Users      Users       `yaml:"Users"`
	Redis      RedisConfig `yaml:"Redis"`
}
type UsersApi struct {
	Port       string `yaml:"port"`
	SECRET_KEY string `yaml:"SECRET_KEY"`
}
type Users struct {
	Port       string `yaml:"port"`
	SECRET_KEY string `yaml:"SECRET_KEY"`
}
type ManagerApi struct {
	Port       string `yaml:"port"`
	SECRET_KEY string `yaml:"SECRET_KEY"`
}

type Database struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DBType   string `yaml:"DBType"`
	DBHost   string `yaml:"DBHost"`
	DBProt   string `yaml:"DBProt"`
	DBName   string `yaml:"DBName"`
}

type RedisConfig struct {
	Host     string `yaml:"Host"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

type Filter struct {
	MinGold    int64 `form:"mingold" json:"mingold" xml:"mingold"  binding:"required"`
	MaxGold    int64 `form:"maxgold" json:"maxgold" xml:"maxgold"  binding:"required"`
	Multiple   int64 `form:"multiple" json:"multiple" xml:"multiple"  binding:"required"`
	Diamond    int64 `form:"diamond" json:"diamond" xml:"diamond"`
	Crazy      int64 `form:"crazy" json:"crazy" xml:"crazy"`
	Cold       int64 `form:"cold" json:"cold" xml:"cold"`
	Precise    int64 `form:"precise" json:"precise" xml:"precise"`
	IgnoreSell bool  `form:"ignore_sell" json:"ignore_sell"`
}

// GetCurrentPath Get Current Path
func GetCurrentPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	return dir, nil
}

func DownloadFile(URL, filepath string) error {
	// 创建HTTP客户端，并设置15秒超时时间
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	// 创建HTTP请求
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}

	// 发送HTTP请求
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unable to download file, status code: %d", resp.StatusCode)
	}

	// 创建文件
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 下载文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// UnzipGZ 解压 .gz 文件
func UnzipGZ(gzFilePath, outputFilePath string) error {
	// 打开要解压的 .gz 文件
	gzFile, err := os.Open(gzFilePath)
	if err != nil {
		return err
	}
	defer gzFile.Close()

	// 创建输出文件
	outFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// 创建 gzip.Reader
	gzReader, err := gzip.NewReader(gzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	// 将解压后的数据复制到输出文件
	_, err = io.Copy(outFile, gzReader)
	if err != nil {
		return err
	}

	log.Println("解压完成")
	return nil
}

func CheckGeoIP(OS, CurrentPath string) {
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	GeoPath := strings.Join([]string{CurrentPath, "GeoIP"}, LinkPathStr)
	if !IsExist(GeoPath) {
		os.MkdirAll(GeoPath, 0755)
	}
	// fmt.Println(GeoPath)
	GeoFileOut := strings.Join([]string{CurrentPath, "GeoIP", "GeoLite2-City.mmdb"}, LinkPathStr)
	GeoFile := strings.Join([]string{CurrentPath, "GeoIP", "GeoLite2-City.mmdb.gz"}, LinkPathStr)
	if !IsExist(GeoFileOut) {
		ProxyUri := []string{
			"",
			"https://github.91chi.fun/",
			"https://ghproxy.com/",
			"https://github.abskoop.workers.dev/",
			"https://gh.api.99988866.xyz/",
		}
		ghuri := "https://cdn.jsdelivr.net/npm/geolite2-city@1.0.0/GeoLite2-City.mmdb.gz"
		for _, item := range ProxyUri {
			uri := strings.Join([]string{item, ghuri}, "")
			err := DownloadFile(uri, GeoFile)
			if err == nil {
				err = UnzipGZ(GeoFile, GeoFileOut)
				if err == nil {
					os.Remove(GeoFile)
					break
				}
			}
		}
	}
}

// CheckConfig check config
func CheckConfig(OS, CurrentPath string) (conf *Config, err error) {
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	ConfigFile := strings.Join([]string{CurrentPath, "config.yaml"}, LinkPathStr)

	var confYaml *Config
	yamlFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}

	if len(confYaml.Port) <= 0 {
		confYaml.Port = "13002"
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.SECRET_KEY) <= 0 {
		secret_key := randSeq(32)
		confYaml.SECRET_KEY = secret_key
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.AdminPWD) <= 0 {
		password := randSeq(12)
		confYaml.AdminPWD = password
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if confYaml.FormMemory == 0 {
		confYaml.FormMemory = 32
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.Database.DBName) <= 0 {
		confYaml.Database.DBName = "acc_Manage"
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.UsersApi.Port) <= 0 {
		confYaml.UsersApi.Port = "13003"
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.UsersApi.SECRET_KEY) <= 0 {
		secret_key := randSeq(32)
		confYaml.UsersApi.SECRET_KEY = secret_key
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.ManagerApi.Port) <= 0 {
		confYaml.ManagerApi.Port = "13005"
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.ManagerApi.SECRET_KEY) <= 0 {
		secret_key := randSeq(32)
		confYaml.ManagerApi.SECRET_KEY = secret_key
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.Users.Port) <= 0 {
		confYaml.Users.Port = "13006"
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	if len(confYaml.Users.SECRET_KEY) <= 0 {
		secret_key := randSeq(32)
		confYaml.Users.SECRET_KEY = secret_key
		config, _ := yaml.Marshal(&confYaml)
		os.WriteFile(ConfigFile, config, 0644)
	}
	return confYaml, nil
}

// CORSMiddleware cors middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// SetConfigMiddleWare set config
func SetConfigMiddleWare(SECRET_KEY, CurrentPath, Users_SECRET_KEY string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("secret_key", SECRET_KEY)
		c.Set("users_secret_key", Users_SECRET_KEY)
		c.Set("current_path", CurrentPath)
		c.Writer.Status()
	}
}

func randSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func IsExist(path string) bool {
	// 判断文件是否存在
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetDateTime() (int64, int64, int64) {
	d := time.Now()
	date := d.Format("2006-01-02")
	//获取当前时区
	loc, _ := time.LoadLocation("Local")

	//日期当天0点时间戳(拼接字符串)
	startDate := date + "_00:00:00"
	startTime, _ := time.ParseInLocation("2006-01-02_15:04:05", startDate, loc)

	//日期当天23时59分时间戳
	endDate := date + "_23:59:59"
	end, _ := time.ParseInLocation("2006-01-02_15:04:05", endDate, loc)

	yesterday := d.AddDate(0, 0, -1)
	yday := yesterday.Format("2006-01-02")
	yDate := yday + "_00:00:00"
	yTime, _ := time.ParseInLocation("2006-01-02_15:04:05", yDate, loc)

	//返回当天0点和23点59分的时间戳
	return startTime.Unix(), end.Unix(), yTime.Unix()
}

func GetDateTimeStr() string {
	d := time.Now()
	date := d.Format("20060102150405")
	return date
}

func GetSqlDateTime(date string) (int64, int64) {
	//获取当前时区
	loc, _ := time.LoadLocation("Local")

	//日期当天0点时间戳(拼接字符串)
	startDate := date + "_00:00:00"
	startTime, _ := time.ParseInLocation("2006-01-02_15:04:05", startDate, loc)

	//日期当天23时59分时间戳
	endDate := date + "_23:59:59"
	end, _ := time.ParseInLocation("2006-01-02_15:04:05", endDate, loc)

	//返回当天0点和23点59分的时间戳
	return startTime.Unix() * 1000, end.Unix() * 1000
}

func MD5(a string) string {
	data := []byte(a)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 字符全部大写
func ConvertToUpperCase(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return unicode.ToUpper(r)
		}
		return r
	}, str)
}

// 判断字符串是否包含特殊字符
func ContainsSpecialCharacters(str string) bool {
	regex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return regex.MatchString(str)
}

func Decimal(num float64) float64 {
	decimal := 2
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	res := strconv.FormatFloat(math.Floor(num*d)/d, 'f', -1, 64)
	floatNum, _ := strconv.ParseFloat(res, 64)
	return floatNum
}

func ReplaceFromThirdChar(str string, length int) string {
	if len(str) <= length {
		return str
	}

	// 将字符串转换为字符切片
	strChars := []rune(str)

	// 从第三位开始替换为 "*"
	for i := length; i < len(strChars); i++ {
		strChars[i] = '*'
	}

	// 将字符切片转换回字符串
	return string(strChars)
}

func ConvertNumber(num int64) string {
	const (
		tenThousand       = 10000
		oneHundredMillion = 100000000
	)

	if num >= oneHundredMillion {
		// 大于等于一亿
		value := float64(num) / float64(oneHundredMillion)
		return strings.Join([]string{strconv.FormatFloat(value, 'f', 1, 64), "亿"}, "")
	} else if num >= tenThousand {
		// 大于等于一万
		value := float64(num) / float64(tenThousand)
		return strings.Join([]string{strconv.FormatFloat(value, 'f', 1, 64), "万"}, "")
	}

	return strconv.FormatInt(num, 10)
}
