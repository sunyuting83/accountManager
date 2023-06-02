package database

import (
	"colaAPI/Manager/utils"
	"database/sql"
	"os"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Eloquent *sql.DB
	sqlDB    *gorm.DB
	DBType   string
)

// InitDB init db
func InitDB(confYaml *utils.Config) {
	GetDB(confYaml)
	DBType = confYaml.Database.DBType
	Eloquent, _ = sqlDB.DB()
	Eloquent.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	Eloquent.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	Eloquent.SetConnMaxLifetime(time.Hour)
	sqlDB.AutoMigrate(&CoinManager{})
}

func GetDB(confYaml *utils.Config) {
	switch confYaml.Database.DBType {
	case "pgsql":
		DNString := strings.Join([]string{"host=", confYaml.Database.DBHost, " user=", confYaml.Database.Username, " password=", confYaml.Database.Password, " dbname=", confYaml.Database.DBName, " port=", confYaml.Database.DBProt, " sslmode=disable TimeZone=Asia/Shanghai"}, "")
		sqlDB, _ = gorm.Open(postgres.New(postgres.Config{
			DSN:                  DNString,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	case "sqlite":
		CurrentPath, _ := utils.GetCurrentPath()
		SqlitePath := MakeSqlitePath(CurrentPath)
		dbPath := strings.Join([]string{SqlitePath, "db"}, "/")
		if !utils.IsExist(dbPath) {
			os.MkdirAll(dbPath, 0755)
		}
		dbName := strings.Join([]string{confYaml.Database.DBName, "db"}, ".")
		dbFile := strings.Join([]string{dbPath, dbName}, "/")
		sqlDB, _ = gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	}
}

func MakeSqlitePath(a string) (d string) {
	b := strings.Split(a, "/")
	l := len(b) - 1
	c := b[0:l]
	d = strings.Join(c, "/")
	return
}