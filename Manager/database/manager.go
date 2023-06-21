package database

type Manager struct {
	ID           uint   `gorm:"primaryKey"`
	UserName     string `gorm:"index"`
	Password     string
	NewStatus    int `gorm:"index"`
	SplitManager SplitManager
	CreatedAt    int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt    int64 `gorm:"autoUpdateTime:milli"`
}

func (manager *Manager) GetCount() (count int64, err error) {
	if err = sqlDB.Model(&manager).Count(&count).Error; err != nil {
		return
	}
	return
}

// Admin List
func GetManagerList(page, Limit int) (manages *[]Manager, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Select("id, user_name, new_status, created_at").
		Preload("SplitManager").
		Order("id desc").
		Limit(Limit).Offset(p).
		Find(&manages, "new_status = ? AND id != 1", 0).Error; err != nil {
		return
	}
	return
}

func CheckanagerID(id uint) (manager *Manager, err error) {
	if err = sqlDB.First(&manager, "id = ? ", id).Error; err != nil {
		return
	}
	return
}
