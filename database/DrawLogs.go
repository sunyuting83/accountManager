package database

type DrawLogs struct {
	ID         uint `gorm:"primaryKey"`
	ProjectsID uint
	DrawUser   string
	LogName    string
	Data       string
	Projects   Projects
	CreatedAt  int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt  int64 `gorm:"autoUpdateTime:milli"`
}

func (draw *DrawLogs) AddDrawLogs() {
	sqlDB.Create(&draw)
}

func (draw *DrawLogs) GetCount(projectid string) (count int64, err error) {
	if err = sqlDB.Model(&draw).Where("projects_id = ?", projectid).Count(&count).Error; err != nil {
		return
	}
	return
}

// Admin List
func GetDrawList(page, Limit int, projectid string) (draw []*DrawLogs, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Preload("Projects").
		Select("ID", "DrawUser", "LogName", "UpdatedAt", "ProjectsID").
		Where("projects_id = ?", projectid).
		Order("id desc").
		Limit(Limit).Offset(p).
		Find(&draw).Error; err != nil {
		return
	}
	return
}

func GetDrawData(id string) (drawed *DrawLogs, err error) {
	if err = sqlDB.
		Preload("Projects").
		Select("Data", "ProjectsID").
		First(&drawed, "id = ? ", id).Error; err != nil {
		return
	}
	return
}
