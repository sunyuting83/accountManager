package database

type Projects struct {
	ID           uint `gorm:"primaryKey"`
	UsersID      uint
	ProjectsName string
	UserName     string
	Password     string
	StatusJSON   string `gorm:"column:status_json;type:longtext"`
	NewStatus    int    `gorm:"index"`
	Accounts     []Accounts
	Filed        []Filed
	Key          string
	Remarks      string `gorm:"column:remarks;type:longtext"`
	AccNumber    int
	ColaAPI      bool
	CreatedAt    int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt    int64 `gorm:"autoUpdateTime:milli"`
}

// Get Count
func (projects *Projects) GetCount(userid int64) (count int64, err error) {
	if err = sqlDB.Model(&projects).Where("users_id = ?", userid).Count(&count).Error; err != nil {
		return
	}
	return
}

// Admin List
func GetProjectsList(userid int64, page, Limit int) (projects *[]Projects, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("users_id = ?", userid).
		Order("projects.id desc").
		Limit(Limit).Offset(p).
		Find(&projects).Error; err != nil {
		return
	}
	return
}

// Check ID
func ProjectsCheckID(id string) (projects *Projects, err error) {
	if err = sqlDB.First(&projects, "id = ?", id).Error; err != nil {
		return
	}
	return
}

// Reset Password
func (projects *Projects) UpdateProjects(id string) {
	sqlDB.Model(&projects).
		Select("UserName", "Password", "AccNumber").
		Where("id = ?", id).
		Updates(&projects)
}
