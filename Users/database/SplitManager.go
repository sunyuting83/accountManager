package database

type SplitManager struct {
	ID        uint `gorm:"primaryKey"`
	ManagerID uint
	Percent   float64
}

func CheckSplitManagerID(id uint) (splitManager *SplitManager, err error) {
	if err = sqlDB.First(&splitManager, SplitManager{ManagerID: id}).Error; err != nil {
		return
	}
	return
}
