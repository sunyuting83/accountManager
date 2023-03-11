package database

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
func GetAccountListUseIn(ProjectsID string, statusList []int) (accounts []Accounts, err error) {
	if err = sqlDB.
		Where("projects_id = ? and new_status in ?", ProjectsID, statusList).
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

// Reset Password
func (account *Accounts) AccountUpStatus(status string) {
	sqlDB.Model(&account).Update("new_status", status)
}

func AccountBatches(accounts []Accounts) {
	sqlDB.Create(&accounts)
}
func AccountInBatches(accounts []Accounts) {
	sqlDB.CreateInBatches(&accounts, 1000)
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
