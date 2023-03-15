package database

import (
	"gorm.io/gorm"
)

type Accounts struct {
	ID            uint `gorm:"primaryKey"`
	ProjectsID    uint
	ComputID      uint
	PhoneNumber   string
	PhonePassword string
	UserName      string
	Password      string
	Cover         string
	NewStatus     int `gorm:"index"`
	TodayGold     int64
	YesterdayGold int64
	Multiple      int64
	Diamond       int
	Crazy         int
	Precise       int
	Cold          int
	Exptime       int64
	Price         float64
	Remarks       string `gorm:"column:remarks;type:longtext"`
	CreatedAt     int64  `gorm:"autoUpdateTime:milli"`
	UpdatedAt     int64  `gorm:"autoUpdateTime:milli"`
}

// Get Count
func (accounts *Accounts) GetCount(ProjectsID, Status string) (count int64, err error) {
	if err = sqlDB.Model(&accounts).Where("projects_id = ? and new_status = ?", ProjectsID, Status).Count(&count).Error; err != nil {
		return
	}
	return
}

// Account List
func (account *Accounts) GetInCount(ProjectsID string, statusList []string) (count int64, err error) {
	if err = sqlDB.
		Model(&account).
		Where("projects_id = ? and new_status IN ?", ProjectsID, statusList).
		Count(&count).Error; err != nil {
		return
	}
	return
}

// Account List
func (account *Accounts) GetInList(ProjectsID string, statusList []string, page, Limit int) (accounts []*Accounts, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("projects_id = ? and new_status IN ?", ProjectsID, statusList).
		Order("updated_at desc").
		Order("updated_at desc").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

// Account List
func GetAccountList(page, Limit int, ProjectsID, Status string) (accounts *[]Accounts, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("projects_id = ? and new_status = ?", ProjectsID, Status).
		Order("updated_at desc").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

// Account List
func GetAccountListUseIn(ProjectsID string, statusList []string) (accounts []Accounts, err error) {
	if err = sqlDB.
		Where("projects_id = ? and new_status IN ?", ProjectsID, statusList).
		Order("updated_at desc").
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

func CheckAccount(projectsid, account string) (accounts *Accounts, err error) {
	if err = sqlDB.First(&accounts, "projects_id = ? and user_name = ? ", projectsid, account).Error; err != nil {
		return
	}
	return
}

// Delete Admin
func (accounts *Accounts) DeleteAll(projectid string, status string) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Where("projects_id = ? and new_status = ?", projectid, status).Delete(&accounts)
}

// Reset Password
func (account *Accounts) AccountUpStatus(status string) {
	sqlDB.Model(&account).Update("new_status", status)
}

// Reset Password
func (account *Accounts) BackTo(projectsID, status, backToStatus string) {
	sqlDB.Model(&account).
		Where("projects_id = ? and new_status = ?", projectsID, status).
		Update("new_status", backToStatus)
}

func AccountBatches(accounts []Accounts) {
	sqlDB.Create(&accounts)
}
func AccountInBatches(accounts []Accounts) {
	sqlDB.CreateInBatches(&accounts, 1000)
}

func (account *Accounts) ExportAccount(projectsID, status string) (accounts []*Accounts, err error) {
	if err = sqlDB.
		Where("projects_id = ? and new_status = ?", projectsID, status).
		Find(&accounts).Error; err != nil {
		return
	}
	return
	// SELECT DISTINCT DATE(updated_at / 1000, 'unixepoch','localtime') FROM accounts WHERE new_status IN (2,3,4,5)
}

func (account *Accounts) PullDataUseIn(IDs []int) {
	sqlDB.
		Model(&account).
		Where("id IN ?", IDs).
		Update("new_status", "108")
}

func (account *Accounts) PullDataUseSQL(SQL string) {
	sqlDB.Exec(SQL)
}

func (account *Accounts) GetDateInCount(projectsID string, statusList []string, starTime, endTime int64) (count int64, err error) {
	if err = sqlDB.
		Model(&account).
		Where("projects_id = ? AND new_status IN ? AND updated_at >= ? AND updated_at <= ?", projectsID, statusList, starTime, endTime).
		Count(&count).Error; err != nil {
		return
	}
	return
}

func GetDateInData(projectsID string, statusList []string, starTime, endTime int64, page, Limit int) (accounts []*Accounts, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("projects_id = ? AND new_status IN ? AND updated_at >= ? AND updated_at <= ?", projectsID, statusList, starTime, endTime).
		Order("updated_at desc").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

func (account *Accounts) GetDatedInCount(projectsID string, starTime, endTime int64) (count int64, err error) {
	if err = sqlDB.
		Model(&account).
		Where("projects_id = ? AND new_status = ? AND updated_at >= ? AND updated_at <= ?", projectsID, "108", starTime, endTime).
		Count(&count).Error; err != nil {
		return
	}
	return
}

func GetDatedInData(projectsID string, starTime, endTime int64, page, Limit int) (accounts []*Accounts, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("projects_id = ? AND new_status = ? AND updated_at >= ? AND updated_at <= ?", projectsID, "108", starTime, endTime).
		Order("updated_at desc").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

func GetDateTimeData(projectsID, statusList, GeType string) (re []string, err error) {
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
	sql := SQLStart + projectsID + " AND new_status IN (" + statusList + ") ORDER BY " + d + " DESC"
	if DBType == "pgsql" {
		sql = SQLStart + projectsID + " AND new_status IN (" + statusList + ") ORDER BY to_char(to_timestamp(" + d + " / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') DESC"
	}
	re, err = RawQueryParseToMap(sqlDB, sql, d)
	return
}

func GetDateTimeDataDraw(projectsID, GeType string) (re []string, err error) {
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
		sql = SQLStart + projectsID + " AND new_status = 108 ORDER BY to_char(to_timestamp(" + d + " / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') DESC"
	}
	re, err = RawQueryParseToMap(sqlDB, sql, d)
	return
}

// makePage make page
func makePage(p, Limit int) int {
	p = p - 1
	if p <= 0 {
		p = 0
	}
	page := p * Limit
	return page
}

// RawQuerySearchAndParseToMap ...
func RawQueryParseToMap(db *gorm.DB, query, date string) ([]string, error) {
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
