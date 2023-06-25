package database

type Order struct {
	ID          uint   `gorm:"primaryKey"`
	OrderCode   string `gorm:"index"`
	Coin        float64
	CoinUsersID uint
	AccountsID  string
	CreatedAt   int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt   int64 `gorm:"autoUpdateTime:milli"`
}

func (order *Order) Insert() (err error) {
	sqlDB.Create(&order)
	return nil
}
