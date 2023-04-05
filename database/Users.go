package database

type Users struct {
	ID        uint `gorm:"primaryKey"`
	ManagerID uint
	UserName  string `gorm:"index"`
	Password  string `gorm:"index"`
	NewStatus int    `gorm:"index"`
	Remarks   string
	Manager   Manager
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

func (user *Users) Insert() (err error) {
	sqlDB.Create(&user)
	return nil
}

// Delete Admin
func (user *Users) DeleteOne(id int64) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Where("id = ?", id).Delete(&user)
}

// Check ID
func UserCheckID(id int64) (user *Users, err error) {
	if err = sqlDB.First(&user, "id = ?", id).Error; err != nil {
		return
	}
	return
}

// Get Count
func (user *Users) GetCount() (count int64, err error) {
	if err = sqlDB.Model(&user).Count(&count).Error; err != nil {
		return
	}
	return
}

// Admin List
func GetUsersList(page, Limit int, manager_id uint) (user *[]Users, err error) {
	p := makePage(page, Limit)
	if manager_id == 1 {
		if err = sqlDB.
			Model(&Users{}).
			Preload("Manager").
			Select("id, user_name, new_status, manager_id, remarks, created_at").
			Order("users.id desc").
			Limit(Limit).Offset(p).
			Find(&user).Error; err != nil {
			return
		}
	} else {
		if err = sqlDB.
			Select("id, user_name, new_status, manager_id, remarks, created_at").
			Where("manager_id = ?", manager_id).
			Order("id desc").
			Limit(Limit).Offset(p).
			Find(&user).Error; err != nil {
			return
		}
	}
	return
}

func GetAllUsersList(manager_id uint) (user *[]Users, err error) {
	if manager_id == 1 {
		if err = sqlDB.
			Model(&Users{}).
			Preload("Manager").
			Order("id desc").
			Find(&user).Error; err != nil {
			return
		}
	} else {
		if err = sqlDB.
			Where("manager_id = ?", manager_id).
			Order("id desc").
			Find(&user).Error; err != nil {
			return
		}
	}
	return
}

// Reset Password
func (user *Users) UserResetPassword(username string) (users Users, err error) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	if err = sqlDB.First(&users, "user_name = ?", username).Error; err != nil {
		return
	}
	if err = sqlDB.Model(&user).Updates(&users).Error; err != nil {
		return
	}
	return
}

// Reset Password
func (user *Users) UserUpStatusAdmin(status int) {
	sqlDB.Model(&user).Update("new_status", status)
}

func (user *Users) SetUserRemarks(remarks string) {
	sqlDB.Model(&user).Update("remarks", remarks)
}

// api

func UserCheckUser(username int, password string) (user *Users, err error) {
	if err = sqlDB.First(&user, "user_name = ? AND new_status = ? AND password = ?", username, "0", password).Error; err != nil {
		return
	}
	return
}
