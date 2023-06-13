package database

type CoinUsers struct {
	ID        uint   `gorm:"primaryKey"`
	UserName  string `gorm:"index"`
	Password  string
	NewStatus int `gorm:"index"`
	Coin      float64
	CreatedAt int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func CheckCoinUserID(id string) (coinuser *CoinUsers, err error) {
	if err = sqlDB.First(&coinuser, "id = ? ", id).Error; err != nil {
		return
	}
	return
}
