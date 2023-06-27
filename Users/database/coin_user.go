package database

import "gorm.io/gorm"

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
	Email         string
	PhoneNumber   string
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

func UserCheckUser(username, password string) (user *CoinUsers, err error) {
	if err = sqlDB.First(&user, "user_name = ? AND new_status = ? AND password = ?", username, "0", password).Error; err != nil {
		return
	}
	return
}

// Check ID
func UserCheckID(id int64) (user *CoinUsers, err error) {
	if err = sqlDB.First(&user, "id = ?", id).Error; err != nil {
		return
	}
	return
}

// Reset Password
func (user *CoinUsers) UserResetPassword(id int64) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Model(&user).Where("id = ?", id).Updates(&user)
}

func UpCoinToCoinUser(id uint, Coin float64) {
	sqlDB.Model(&CoinUsers{}).
		Where("id = ?", id).
		Update("coin", gorm.Expr("coin - ?", Coin))
}
