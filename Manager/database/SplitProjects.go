package database

type SplitProjects struct {
	ID      uint `gorm:"primaryKey"`
	Percent float64
	Manager string
}

func SetSplitProjects(Percent float64, ManagerID string) error {
	var splitProjects *SplitProjects
	result := sqlDB.First(&splitProjects, SplitManager{ID: 1})
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}
	if result.RowsAffected == 0 && result.Error.Error() != "record not found" {
		// 更新字段
		result = sqlDB.UpdateColumns(&SplitProjects{Percent: Percent, Manager: ManagerID})
		if result.Error != nil {
			return result.Error
		}
	} else {
		// fmt.Println("where")
		result = sqlDB.UpdateColumns(&SplitProjects{ID: 1, Percent: Percent, Manager: ManagerID})
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
