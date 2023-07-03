package database

type Transfers struct {
	ID              uint `gorm:"primaryKey"`
	CoinUsersID     uint
	FormCoinUsersID uint `gorm:"foreignKey:CoinUsersID"`
	Coin            float64
	CreatedAt       int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt       int64 `gorm:"autoUpdateTime:milli"`
}
