package database

import "gorm.io/gorm"

type Users struct {
	ID        uint `gorm:"primaryKey"`
	ManagerID uint
	Coin      float64
}

func UpCoinToUsers(Coin float64, id uint) {
	sqlDB.Model(&Users{}).
		Where("id = ?", id).
		Update("coin", gorm.Expr("coin + ?", Coin))
}
