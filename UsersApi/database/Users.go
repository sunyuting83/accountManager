package database

type Users struct {
	ID        uint   `gorm:"primaryKey"`
	UserName  string `gorm:"index"`
	Password  string `gorm:"index"`
	NewStatus int    `gorm:"index"`
	Remarks   string
	Projects  []Projects
	CreatedAt int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func UserCheckUserName(username string) (user *Users, err error) {
	if err = sqlDB.First(&user, "user_name = ? ", username).Error; err != nil {
		return
	}
	return
}

// Check ID
func UserCheckID(id int64) (user *Users, err error) {
	if err = sqlDB.First(&user, "id = ?", id).Error; err != nil {
		return
	}
	return
}

// Reset Password
func (user *Users) UserResetPassword(id int64) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Model(&user).Where("id = ?", id).Updates(&user)
}

// api

func UserCheckUser(username, password string) (user *Users, err error) {
	if err = sqlDB.First(&user, "user_name = ? AND new_status = ? AND password = ?", username, "0", password).Error; err != nil {
		return
	}
	return
}
