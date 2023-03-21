package database

type Filed struct {
	ID         uint `gorm:"primaryKey"`
	ProjectsID uint
	FiledName  string
	Data       string
	CreatedAt  int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt  int64 `gorm:"autoUpdateTime:milli"`
}

func GetFiledList(ProjectsID string) (filed []*Filed, err error) {
	if err = sqlDB.
		Select("filed_name").
		Where("projects_id = ?", ProjectsID).
		Order("created_at DESC").
		Find(&filed).Error; err != nil {
		return
	}
	return
}

func GetFirstFiled(ProjectsID string) (filed *Filed, err error) {
	// fmt.Println(status)
	if err = sqlDB.
		Where("projects_id = ?", ProjectsID).
		Order("created_at DESC").
		First(&filed).Error; err != nil {
		return
	}
	return
}
func GetOneFiled(ProjectsID, FiledName string) (filed *Filed, err error) {
	// fmt.Println(status)
	if err = sqlDB.
		Where("projects_id = ? AND filed_name = ?", ProjectsID, FiledName).
		First(&filed).Error; err != nil {
		return
	}
	return
}
