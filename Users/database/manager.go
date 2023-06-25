package database

type Manager struct {
	ID   uint `gorm:"primaryKey"`
	Coin float64
}

func UpCoinToManager(Coin float64, id uint) {
	sqlDB.Model(&Manager{}).
		Where("id = ?", id).
		UpdateColumns(Manager{Coin: Coin})
}
