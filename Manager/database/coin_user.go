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
	Email         string
	PhoneNumber   string
	WalletAddress string `gorm:"index"`
	CreatedAt     int64  `gorm:"autoUpdateTime:milli"`
	UpdatedAt     int64  `gorm:"autoUpdateTime:milli"`
}

func CheckCoinUserID(id string) (coinuser *CoinUsers, err error) {
	if err = sqlDB.First(&coinuser, "id = ? ", id).Error; err != nil {
		return
	}
	return
}

// Check ID
func UserCheckID(id int64) (coinuser *CoinUsers, err error) {
	if err = sqlDB.First(&coinuser, "id = ?", id).Error; err != nil {
		return
	}
	return
}

func (coinuser *CoinUsers) UpCoinToCoinUsers() {
	sqlDB.Save(&coinuser)
}

func UserCheckUserName(username string) (coinuser *CoinUsers, err error) {
	if err = sqlDB.First(&coinuser, "user_name = ? ", username).Error; err != nil {
		return
	}
	return
}

func UserCheckUserKey(key string) (coinuser *CoinUsers, err error) {
	if err = sqlDB.First(&coinuser, "wallet_address = ? ", key).Error; err != nil {
		return
	}
	return
}

func (coinuser *CoinUsers) Insert() (err error) {
	sqlDB.Create(&coinuser)
	return nil
}

// Delete Admin
func (coinuser *CoinUsers) DeleteOne(id int64) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Where("id = ?", id).Delete(&coinuser)
}

// Get Count
func (coinuser *CoinUsers) GetCount() (count int64, err error) {
	if err = sqlDB.Model(&coinuser).Count(&count).Error; err != nil {
		return
	}
	return
}

// Admin List
func GetUsersList(page, Limit int) (coinuser *[]CoinUsers, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Model(&CoinUsers{}).
		Select("id, user_name, new_status, parent_id, coin, local_address, wallet_address, created_at").
		Order("id desc").
		Limit(Limit).Offset(p).
		Find(&coinuser).Error; err != nil {
		return
	}
	return
}

func GetAllUsersList() (coinuser *[]CoinUsers, err error) {
	if err = sqlDB.
		Model(&CoinUsers{}).
		Order("id desc").
		Find(&coinuser).Error; err != nil {
		return
	}
	return
}

// Reset Password
func (coinuser *CoinUsers) UserResetPassword(username string) (coinusers CoinUsers, err error) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	if err = sqlDB.First(&coinusers, "user_name = ?", username).Error; err != nil {
		return
	}
	if err = sqlDB.Model(&coinuser).Updates(&coinusers).Error; err != nil {
		return
	}
	return
}

// Reset Password
func (coinuser *CoinUsers) UpUserStatus(status int) {
	sqlDB.Model(&coinuser).Update("new_status", status)
}
