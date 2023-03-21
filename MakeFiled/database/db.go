package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Eloquent *sql.DB
	sqlDB    *gorm.DB
	DBType   string
)

type Config struct {
	FiledSize int      `yaml:"FiledSize"`
	Database  Database `yaml:"Database"`
}

type Database struct {
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DBType   string `yaml:"DBType"`
	DBHost   string `yaml:"DBHost"`
	DBProt   string `yaml:"DBProt"`
	DBName   string `yaml:"DBName"`
}

// InitDB init db
func InitDB() {

	confYaml, err := CheckConfig()
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(10) * time.Second)
		os.Exit(0)
	}
	GetDB(confYaml)
	DBType = confYaml.Database.DBType
	Eloquent, _ = sqlDB.DB()
	Eloquent.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	Eloquent.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	Eloquent.SetConnMaxLifetime(time.Hour)

}

func GetDB(confYaml *Config) {
	switch confYaml.Database.DBType {
	case "pgsql":
		DNString := strings.Join([]string{"host=", confYaml.Database.DBHost, " user=", confYaml.Database.Username, " password=", confYaml.Database.Password, " dbname=", confYaml.Database.DBName, " port=", confYaml.Database.DBProt, " sslmode=disable TimeZone=Asia/Shanghai"}, "")
		sqlDB, _ = gorm.Open(postgres.New(postgres.Config{
			DSN:                  DNString,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	case "mysql":
		DNString := strings.Join([]string{confYaml.Database.Username, ":", confYaml.Database.Password, "@tcp(", confYaml.Database.DBHost, ":", confYaml.Database.DBProt, ")/", confYaml.Database.DBName, "?charset=utf8&parseTime=True&loc=Local"}, "")
		sqlDB, _ = gorm.Open(mysql.New(mysql.Config{
			DSN:                       DNString, // DSN data source name
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}), &gorm.Config{})
	case "sqlite":
		CurrentPath, _ := GetCurrentPath()
		dbPath := strings.Join([]string{CurrentPath, "db"}, "/")
		if !IsExist(dbPath) {
			os.MkdirAll(dbPath, 0755)
		}
		dbName := strings.Join([]string{confYaml.Database.DBName, "db"}, ".")
		dbFile := strings.Join([]string{dbPath, dbName}, "/")
		sqlDB, _ = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	}
}

// CheckConfig check config
func CheckConfig() (conf *Config, err error) {

	OS := runtime.GOOS
	CurrentPath, _ := GetCurrentPath()
	ConfigFilePath := MakeSqlitePath(CurrentPath)
	LinkPathStr := "/"
	if OS == "windows" {
		LinkPathStr = "\\"
	}
	ConfigFile := strings.Join([]string{ConfigFilePath, "config.yaml"}, LinkPathStr)

	var confYaml *Config
	yamlFile, err := os.ReadFile(ConfigFile)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	err = yaml.Unmarshal(yamlFile, &confYaml)
	if err != nil {
		return confYaml, errors.New("读取配置文件出错\n10秒后程序自动关闭")
	}
	return confYaml, nil
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

func MakeSqlitePath(a string) (d string) {
	b := strings.Split(a, "/")
	l := len(b) - 1
	c := b[0:l]
	d = strings.Join(c, "/")
	return
}

func IsExist(path string) bool {
	// 判断文件是否存在
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func GetDateTime() int64 {
	confYaml, _ := CheckConfig()
	d := time.Now()
	//获取当前时区
	loc, _ := time.LoadLocation("UTC")

	yesterday := d.AddDate(0, 0, -confYaml.FiledSize)
	yday := yesterday.Format("2006-01-02")
	yDate := yday + "_00:00:00"
	yTime, _ := time.ParseInLocation("2006-01-02_15:04:05", yDate, loc)

	//返回当天0点和23点59分的时间戳
	return yTime.Unix()
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
