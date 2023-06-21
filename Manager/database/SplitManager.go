package database

type SplitManager struct {
	ID        uint `gorm:"primaryKey"`
	ManagerID uint
	Percent   float64
}

func SetSplitManager(id uint, Percent float64) error {
	var splitManager *SplitManager
	result := sqlDB.First(&splitManager, SplitManager{ManagerID: id})
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}
	if result.RowsAffected == 0 && result.Error.Error() != "record not found" {
		// 更新字段
		result = sqlDB.UpdateColumns(SplitManager{Percent: Percent})
		if result.Error != nil {
			return result.Error
		}
	} else {
		splitManager.ManagerID = id
		splitManager.Percent = Percent
		result = sqlDB.Save(&splitManager)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
