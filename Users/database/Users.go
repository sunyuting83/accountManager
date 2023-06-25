package database

type Users struct {
	ID        uint `gorm:"primaryKey"`
	ManagerID uint
	Coin      float64
}

func UpCoinToUsers(Coin float64, id uint) {
	sqlDB.Model(&Users{}).
		Where("id = ?", id).
		UpdateColumns(Users{Coin: Coin})
}
