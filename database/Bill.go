package database

type Bill struct {
	ID              uint `gorm:"primaryKey"`
	CoinUsersID     *uint
	FormCoinUsersID *uint `gorm:"foreignKey:CoinUsersID"`
	OrderID         *uint
	Coin            float64
	NewStatus       int
	Months          string `gorm:"index"`
	CreatedAt       int64  `gorm:"autoUpdateTime:milli"`
	UpdatedAt       int64  `gorm:"autoUpdateTime:milli"`
}
