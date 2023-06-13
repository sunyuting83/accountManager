package database

type CoinUsers struct {
	ID            uint   `gorm:"primaryKey"`
	UserName      string `gorm:"index"`
	Password      string
	NewStatus     int `gorm:"index"`
	Coin          float64
	CoinUsers     []CoinUsers
	IPAddress     string
	LocalAddress  string
	WalletAddress string `gorm:"index"`
	CreatedAt     int64  `gorm:"autoUpdateTime:milli"`
	UpdatedAt     int64  `gorm:"autoUpdateTime:milli"`
}

func (users *CoinUsers) Insert() (err error) {
	sqlDB.Create(&users)
	return nil
}

func CheckUserName(username string) (users *CoinUsers, err error) {
	if err = sqlDB.First(&users, "user_name = ? ", username).Error; err != nil {
		return
	}
	return
}
