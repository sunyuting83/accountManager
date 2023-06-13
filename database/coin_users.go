package database

type CoinUsers struct {
	ID           uint   `gorm:"primaryKey"`
	UserName     string `gorm:"index"`
	Password     string
	NewStatus    int `gorm:"index"`
	Coin         float64
	IPAddress    string
	LocalAddress string
	CreatedAt    int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt    int64 `gorm:"autoUpdateTime:milli"`
}
