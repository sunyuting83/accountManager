package database

import "gorm.io/gorm"

type Manager struct {
	ID   uint `gorm:"primaryKey"`
	Coin float64
}

func UpCoinToManager(Coin float64, id uint) {
	sqlDB.Model(&Manager{}).
		Where("id = ?", id).
		Update("coin", gorm.Expr("coin + ?", Coin))
}
