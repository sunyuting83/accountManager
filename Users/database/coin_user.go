package database

type CoinUsers struct {
	ID            uint   `gorm:"primaryKey"`
	ParentID      *uint  `gorm:"foreignKey:ParentID"`
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

func CheckUserKey(key string) (users *CoinUsers, err error) {
	if err = sqlDB.First(&users, "wallet_address = ? ", key).Error; err != nil {
		return
	}
	return
}
