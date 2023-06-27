package database

type Order struct {
	ID          uint   `gorm:"primaryKey"`
	OrderCode   string `gorm:"index"`
	Coin        float64
	CoinUsersID uint
	Accounts    []Accounts
	CreatedAt   int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt   int64 `gorm:"autoUpdateTime:milli"`
}
