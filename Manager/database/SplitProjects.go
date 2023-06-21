package database

type SplitProjects struct {
	ID      uint `gorm:"primaryKey"`
	Percent float64
}

func SetSplitProjects(Percent float64) error {
	var splitProjects *SplitProjects
	result := sqlDB.First(&splitProjects, SplitManager{ID: 1})
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}
	if result.RowsAffected == 0 && result.Error.Error() != "record not found" {
		// 更新字段
		result = sqlDB.UpdateColumns(SplitProjects{Percent: Percent})
		if result.Error != nil {
			return result.Error
		}
	} else {
		splitProjects.ID = 1
		splitProjects.Percent = Percent
		result = sqlDB.Save(&splitProjects)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
