package database

import "fmt"

type SplitManager struct {
	ID         uint `gorm:"primaryKey"`
	ManagerID  uint
	Proportion float64
}

func (splitManager *SplitManager) Insert() (err error) {
	sqlDB.Save(&splitManager)
	return nil
}

func SetSplitManager(id uint, Proportion float64) error {
	var splitManager *SplitManager
	result := sqlDB.First(&splitManager, SplitManager{ManagerID: id})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		fmt.Println("here")
		// 更新字段
		result = sqlDB.UpdateColumns(SplitManager{Proportion: Proportion})
		if result.Error != nil {
			return result.Error
		}
	} else {
		splitManager.ManagerID = id
		splitManager.Proportion = Proportion
		result = sqlDB.Save(&splitManager)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
