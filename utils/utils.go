package utils

import (
	"colaAPI/UsersApi/utils"
	BadgerDB "colaAPI/badger"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	Redis      RedisConfig `yaml:"Redis"`
}

type UsersApi struct {
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

// GetCurrentPath Get Current Path
func GetCurrentPath() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	return dir, nil
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
		confYaml.AdminPWD = "admin888"
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
func SetConfigMiddleWare(SECRET_KEY, CurrentPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("secret_key", SECRET_KEY)
		c.Set("current_path", CurrentPath)
		c.Writer.Status()
	}
}

func GetTokenUserData(c *gin.Context) (result *CacheToken) {

	token := c.GetHeader("Authorization")

	secret_key, _ := c.Get("secret_key")
	SECRET_KEY := secret_key.(string)
	token = token[7:]
	AEStoken, err := utils.DecryptByAes(token, []byte(SECRET_KEY))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't token",
		})
		return
	}
	Token, err := BadgerDB.GetToken(AEStoken)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	json.Unmarshal(Token, &result)
	return
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

func MD5(a string) string {
	data := []byte(a)
	md5Ctx := md5.New()
	md5Ctx.Write(data)
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
