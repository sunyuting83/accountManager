package database

import (
	"colaAPI/UsersApi/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Accounts struct {
	ID            uint `gorm:"primaryKey"`
	ProjectsID    uint
	GameID        *uint
	ComputID      uint
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

// Get Count
func (accounts *Accounts) GetCount(ProjectsID, Status string) (count int64, err error) {
	if err = sqlDB.Model(&accounts).Where("projects_id = ? and new_status = ?", ProjectsID, Status).Count(&count).Error; err != nil {
		return
	}
	return
}
func (accounts *Accounts) AddAccount() {
	sqlDB.Create(&accounts)
}

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
		Order("today_gold DESC").
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

func CheckOneAccount(projectsid, account string) (accounts *Accounts, err error) {
	if err = sqlDB.First(&accounts, "projects_id = ? and user_name = ? and new_status != ?", projectsid, account, "108").Error; err != nil {
		return
	}
	return
}

func CheckAccountForStatus(projectsid, account, status string) (accounts *Accounts, err error) {
	if err = sqlDB.First(&accounts, "projects_id = ? and user_name = ? and new_status = ?", projectsid, account, status).Error; err != nil {
		return
	}
	return
}

// Delete Admin
func (accounts *Accounts) DeleteAll(projectid string, status string) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Where("projects_id = ? and new_status = ?", projectid, status).Delete(&accounts)
}

// Delete Admin
func (accounts *Accounts) DeleteOne(projectid string, account string) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Where("projects_id = ? and user_name = ? and new_status != ?", projectid, account, "108").Delete(&accounts)
}

// Reset Password
func (account *Accounts) AccountUpStatus(status string) {
	// fmt.Println(status)
	sqlDB.Model(&account).Update("new_status", status)
}

// Reset Password
func (account *Accounts) AccountUpComput(comput uint) {
	sqlDB.Model(&account).Update("comput_id", comput)
}

// Reset Password
func (account *Accounts) BackTo(projectsID, status string, backToStatus int) {
	sqlDB.Model(&account).
		Select("comput_id", "new_status", "updated_at").
		Where("projects_id = ? and new_status = ?", projectsID, status).
		Updates(Accounts{ComputID: uint(0), NewStatus: backToStatus})
}

func (account *Accounts) UpdataOneAccount(projectsID, username string, accounts map[string]interface{}) {
	sqlDB.Model(&account).
		Omit("created_at").
		Where("projects_id = ? and user_name = ?", projectsID, username).
		Updates(accounts)
}

// Reset Password
func (account *Accounts) BackToAcc(projectsID, status string, backToStatus int, comput uint) {
	sqlDB.Model(&account).
		Where("projects_id = ? and new_status = ? and comput_id = ?", projectsID, status, comput).
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

func PullDataUseIn(IDs []int) (accounts []*Accounts, err error) {
	sqlDB.
		Model(&accounts).
		Clauses(clause.Returning{}).
		Where("id IN ?", IDs).
		Updates(Accounts{NewStatus: 108, SellStatus: 2})
	return
}
func SetSellUseIn(IDs []int) (accounts []*Accounts, err error) {
	sqlDB.
		Model(&accounts).
		Clauses(clause.Returning{}).
		Where("id IN ?", IDs).
		UpdateColumns(Accounts{SellStatus: 1})
	return
}

func MinGold(MinGold int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if MinGold > 0 {
			return db.Where("today_gold >= ?", MinGold)
		}
		return db.Where("")
	}
}
func MaxGold(MaxGold int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if MaxGold > 0 {
			return db.Where("today_gold <= ?", MaxGold)
		}
		return db.Where("")
	}
}
func Multiple(Multiple int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Multiple > 0 {
			return db.Where("multiple >= ?", Multiple)
		}
		return db.Where("")
	}
}
func Diamond(Diamond int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Diamond > 0 {
			return db.Where("diamond >= ?", Diamond)
		}
		return db.Where("")
	}
}
func Crazy(Crazy int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Crazy > 0 {
			return db.Where("crazy >= ?", Crazy)
		}
		return db.Where("")
	}
}
func Cold(Cold int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Cold > 0 {
			return db.Where("cold >= ?", Cold)
		}
		return db.Where("")
	}
}
func Precise(Precise int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Precise > 0 {
			return db.Where("precise >= ?", Precise)
		}
		return db.Where("")
	}
}
func HasStatus(hasStatus []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("new_status IN (?)", hasStatus)
	}
}
func IngoreSell(Ignore bool) func(db *gorm.DB) *gorm.DB {
	if Ignore {
		return func(db *gorm.DB) *gorm.DB {
			return db.Where("sell_status = ?", 0)
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("")
	}
}

func GetDataUseScopes(filter utils.Filter, hasStatus []string, projectsID string) (accounts []*Accounts, err error) {
	if err = sqlDB.
		Where("projects_id = ?", projectsID).
		Scopes(HasStatus(hasStatus)).
		Scopes(MinGold(filter.MinGold)).
		Scopes(MaxGold(filter.MaxGold)).
		Scopes(Multiple(filter.Multiple)).
		Scopes(Diamond(filter.Diamond)).
		Scopes(Crazy(filter.Crazy)).
		Scopes(Cold(filter.Cold)).
		Scopes(Precise(filter.Precise)).
		Order("today_gold DESC").
		Find(&accounts).Error; err != nil {
		return
	}
	return
}
func GetDataUseScopesB(filter utils.Filter, hasStatus []string, projectsID string, Ignore bool) (accounts []*Accounts, err error) {
	if err = sqlDB.
		Where("projects_id = ?", projectsID).
		Scopes(IngoreSell(Ignore)).
		Scopes(HasStatus(hasStatus)).
		Scopes(MinGold(filter.MinGold)).
		Scopes(MaxGold(filter.MaxGold)).
		Scopes(Multiple(filter.Multiple)).
		Scopes(Diamond(filter.Diamond)).
		Scopes(Crazy(filter.Crazy)).
		Scopes(Cold(filter.Cold)).
		Scopes(Precise(filter.Precise)).
		Order("today_gold DESC").
		Find(&accounts).Error; err != nil {
		return
	}
	return
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
		Order("today_gold desc").
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
		Order("today_gold desc").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

func ExportAccountDrawed(projectsID string, starTime, endTime int64) (accounts []*Accounts, err error) {
	if err = sqlDB.
		Where("projects_id = ? AND new_status = ? AND updated_at >= ? AND updated_at <= ?", projectsID, "108", starTime, endTime).
		Order("today_gold desc").
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
	sql := SQLStart + projectsID + " AND new_status = 108 ORDER BY " + d + " DESC"
	// fmt.Println(sql)

	if DBType == "pgsql" {
		sql = SQLStart + projectsID + " AND new_status = 108 ORDER BY to_char(to_timestamp(" + d + " / 1000) AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD') DESC"
	}
	re, err = RawQueryParseToMap(sqlDB, sql, d)
	return
}

func GetOneAccount(ProjectsID, status string) (accounts *Accounts, err error) {
	// fmt.Println(status)
	if err = sqlDB.
		Where("projects_id = ? and new_status = ?", ProjectsID, status).
		First(&accounts).Error; err != nil {
		return
	}
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
