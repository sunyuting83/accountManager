package database

import (
	"gorm.io/gorm"
)

type Accounts struct {
	ID            uint `gorm:"primaryKey"`
	ProjectsID    uint
	GameID        uint
	ComputID      uint
	OrderID       uint
	PhoneNumber   string
	PhonePassword string
	UserName      string
	Password      string
	Cover         string
	NewStatus     int `gorm:"index"`
	SellStatus    int `gorm:"index;default:0"`
	TodayGold     int64
	YesterdayGold int64
	Multiple      int64
	Diamond       int
	Crazy         int
	Precise       int
	Cold          int
	Exptime       int64
	Price         float64
	Remarks       string
	CreatedAt     int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt     int64 `gorm:"autoUpdateTime:milli"`
}

func GetDateHasProject(GeType, checkDate string) (re []int64, err error) {

	SQLStart := "SELECT DISTINCT projects_id FROM accounts WHERE to_char(to_timestamp(updated_at / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') <= to_char(to_timestamp(" + checkDate + " / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') "
	if DBType == "pgsql" {
		SQLStart = "SELECT DISTINCT projects_id FROM accounts WHERE to_char(to_timestamp(updated_at / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') <= to_char(to_timestamp(" + checkDate + " / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') "
	}
	if DBType == "mysql" {
		SQLStart = "SELECT DISTINCT projects_id FROM accounts WHERE updated_at <= "
	}
	sql := SQLStart + " AND new_status = 108"
	// fmt.Println(sql)
	re, err = RawQueryParseToMap(sqlDB, sql)
	return
}

func GetDateTimeDataDraw(projectsID, GeType, checkDate string) (re []string, err error) {
	var d string = "updated_at"
	if GeType == "0" {
		d = "created_at"
	}
	SQLStart := "SELECT DISTINCT DATE(" + d + " / 1000 ,'unixepoch','localtime') FROM accounts WHERE projects_id = "
	if DBType == "pgsql" {
		SQLStart = "SELECT DISTINCT to_char(to_timestamp(" + d + " / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') FROM accounts WHERE projects_id = "
	}
	if DBType == "mysql" {
		SQLStart = "SELECT DISTINCT DATE_FORMAT(from_unixtime(" + d + ` / 1000) ,'%Y-%m-%d') FROM accounts WHERE projects_id = `
	}
	sql := SQLStart + projectsID + " AND new_status = 108 ORDER BY " + d + " DESC"
	// fmt.Println(sql)

	if DBType == "pgsql" {
		sql = SQLStart + projectsID + " AND new_status = 108 AND to_char(to_timestamp(updated_at / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') <= to_char(to_timestamp(" + checkDate + " / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD')  ORDER BY to_char(to_timestamp(" + d + " / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') DESC"
	}
	re, err = RawQueryParseToMapa(sqlDB, sql, d)
	return
}

// Delete Admin
func DeleteAtDate(idList []int) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Delete(&Accounts{}, idList)
}

func GetDateInData(projectsID string, starTime, endTime int64) (accounts []*Accounts, err error) {
	if err = sqlDB.
		Where("projects_id = ? AND new_status = ? AND updated_at >= ? AND updated_at <= ?", projectsID, "108", starTime, endTime).
		Order("today_gold desc").
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

// RawQuerySearchAndParseToMap ...
func RawQueryParseToMap(db *gorm.DB, query string) ([]int64, error) {
	//Use raw query
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	//取得搜尋回來的資料所擁有的column
	columns, er := rows.Columns()
	if er != nil {
		return nil, er
	}
	columnLength := len(columns)

	//make一個臨時儲存的地方，並賦予指標
	cache := make([]interface{}, columnLength)
	for index := range cache {
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]interface{}
	for rows.Next() {
		rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //column可能有許多種data type，因此在這取出時不指定型別，否則會轉換錯誤，且在這取出時為uint8(btye array)格式
		}

		list = append(list, item)
	}
	var l []int64
	for index := range list {
		for _, column := range columns {
			l = append(l, list[index][column].(int64))
		}
	}
	return l, nil
}

// RawQuerySearchAndParseToMap ...
func RawQueryParseToMapa(db *gorm.DB, query, date string) ([]string, error) {
	//Use raw query
	rows, err := db.Raw(query).Rows()
	if err != nil {
		return nil, err
	}

	//取得搜尋回來的資料所擁有的column
	columns, er := rows.Columns()
	if er != nil {
		return nil, er
	}
	columnLength := len(columns)

	//make一個臨時儲存的地方，並賦予指標
	cache := make([]interface{}, columnLength)
	for index := range cache {
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]interface{}
	for rows.Next() {
		rows.Scan(cache...)

		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{}) //column可能有許多種data type，因此在這取出時不指定型別，否則會轉換錯誤，且在這取出時為uint8(btye array)格式
		}

		list = append(list, item)
	}
	var l []string
	//將byte array轉換為字串
	DateFunction := "DATE(" + date + " / 1000 ,'unixepoch','localtime')"
	if DBType == "pgsql" {
		DateFunction = "to_char"
	}
	if DBType == "mysql" {
		DateFunction = "DATE_FORMAT(from_unixtime(" + date + ` / 1000) ,'%Y-%m-%d')`
	}
	for index := range list {
		for _, column := range columns {
			if column == DateFunction {
				// list[index][column] = list[index][column].(string)
				// if DBType == "mysql" {
				// 	l = append(l, string(list[index][column].([]uint8)))
				// } else {
				// 	l = append(l, list[index][column].(string))
				// }
				l = append(l, list[index][column].(string))
			}
		}
	}
	return l, nil
}
