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

func (order *Order) Insert() (id uint, err error) {
	result := sqlDB.Create(&order)
	if result.Error != nil {
		return 0, result.Error
	}
	return order.ID, nil
}
