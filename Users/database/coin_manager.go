package database

import "gorm.io/gorm"

type CoinManager struct {
	ID        uint   `gorm:"primaryKey"`
	UserName  string `gorm:"index"`
	Password  string
	NewStatus int `gorm:"index"`
	Coin      float64
	CreatedAt int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func UpCoinToCoinManager(Coin float64, id []string) {
	for _, item := range id {
		sqlDB.Model(&CoinManager{}).
			Where("id = ?", item).
			Update("coin", gorm.Expr("coin + ?", Coin))
	}
}
