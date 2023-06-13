package database

type CoinUsers struct {
	ID            uint   `gorm:"primaryKey"`
	ParentID      uint   `gorm:"foreignKey:ParentID"`
	UserName      string `gorm:"index"`
	Password      string
	NewStatus     int `gorm:"index"`
	Coin          float64
	Children      []CoinUsers `gorm:"foreignKey:ParentID"`
	IPAddress     string
	LocalAddress  string
	WalletAddress string `gorm:"index"`
	CreatedAt     int64  `gorm:"autoUpdateTime:milli"`
	UpdatedAt     int64  `gorm:"autoUpdateTime:milli"`
}
