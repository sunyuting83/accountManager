package database

import "fmt"

type CoinManager struct {
	ID        uint   `gorm:"primaryKey"`
	UserName  string `gorm:"index"`
	Password  string
	NewStatus int `gorm:"index"`
	Coin      float64
	CreatedAt int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func (manager *CoinManager) Insert() (err error) {
	sqlDB.Create(&manager)
	return nil
}

// Check ID
func CoinManagerCheckID(id int64) (manager *CoinManager, err error) {
	if err = sqlDB.First(&manager, "id = ?", id).Error; err != nil {
		return
	}
	return
}

func CheckAdminLogin(username, password string) (manager *CoinManager, err error) {
	if err = sqlDB.First(&manager, "user_name = ? AND new_status = ? AND password = ?", username, "0", password).Error; err != nil {
		return
	}
	return
}

func CheckUserName(username string) (manager *CoinManager, err error) {
	if err = sqlDB.First(&manager, "user_name = ? ", username).Error; err != nil {
		return
	}
	return
}

func CheckUserID(id uint) (manager *CoinManager, err error) {
	if err = sqlDB.First(&manager, "id = ? ", id).Error; err != nil {
		return
	}
	return
}

// Get Count
func (manager *CoinManager) GetCount() (count int64, err error) {
	if err = sqlDB.Model(&manager).Count(&count).Error; err != nil {
		return
	}
	return
}

// Check ID
func CheckID(id int64) (manager *CoinManager, err error) {
	if err = sqlDB.First(&manager, "id = ?", id).Error; err != nil {
		return
	}
	return
}

// Delete Admin
func (manager *CoinManager) DeleteOne(id int64) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Where("id = ?", id).Delete(&manager)
}

// Admin List
func GetAdminList(page, Limit int) (manages *[]CoinManager, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Select("id, user_name, new_status, created_at").
		Order("id desc").
		Limit(Limit).Offset(p).
		Find(&manages).Error; err != nil {
		return
	}
	return
}

// Admin List
func GetHasUsersID(manager_id uint) (manages *CoinManager, err error) {
	if err = sqlDB.
		Where(&CoinManager{ID: manager_id}).Preload("Users").Find(&manages).Error; err != nil {
		return
	}
	return
}

// Reset Password
func (manager *CoinManager) ResetPassword(username string) (manage CoinManager, err error) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	if err = sqlDB.First(&manage, "user_name = ?", username).Error; err != nil {
		return
	}
	fmt.Println(manager)
	if err = sqlDB.Model(&manage).Updates(&manager).Error; err != nil {
		return
	}
	return
}

// Reset Password
func (manager *CoinManager) UpStatusAdmin(status int) {
	sqlDB.Model(&manager).Update("new_status", status)
}

// makePage make page
func makePage(p, Limit int) int {
	p = p - 1
	if p <= 0 {
		p = 0
	}
	page := p * Limit
	return page
}
