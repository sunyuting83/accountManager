package database

import "gorm.io/gorm/clause"

type Projects struct {
	ID           uint `gorm:"primaryKey"`
	UsersID      uint
	GamesID      uint
	ProjectsName string
	UserName     string
	Password     string
	StatusJSON   string
	NewStatus    int `gorm:"index"`
	Accounts     []Accounts
	Filed        []Filed
	Key          string
	Remarks      string
	AccNumber    int
	ColaAPI      bool
	Users        Users
	CreatedAt    int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt    int64 `gorm:"autoUpdateTime:milli"`
}

func (projects *Projects) Insert() (err error) {
	sqlDB.Create(&projects)
	return nil
}

// Get Count
func (projects *Projects) GetCount() (count int64, err error) {
	if err = sqlDB.Model(&projects).Count(&count).Error; err != nil {
		return
	}
	return
}
func (projects *Projects) GetCountWithIn(UsersIDs []int) (count int64, err error) {
	if err = sqlDB.
		Model(&projects).
		Where("users_id IN ?", UsersIDs).
		Count(&count).Error; err != nil {
		return
	}
	return
}

// Admin List
func GetProjectsList(page, Limit int) (projects *[]Projects, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Model(&Projects{}).
		Preload("Users").
		Order("projects.id desc").
		Limit(Limit).Offset(p).
		Find(&projects).Error; err != nil {
		return
	}
	return
}

func GetProjectsListWithIn(page, Limit int, UsersIDs []int) (projects *[]Projects, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Model(&Projects{}).
		Preload("Users").
		Where("users_id IN ?", UsersIDs).
		Order("projects.id desc").
		Limit(Limit).Offset(p).
		Find(&projects).Error; err != nil {
		return
	}
	return
}

// Check ID
func ProjectsCheckID(id int64) (projects *Projects, err error) {
	if err = sqlDB.First(&projects, "id = ?", id).Error; err != nil {
		return
	}
	return
}

// Delete Admin
func (projects *Projects) DeleteOne(id int64) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Select(clause.Associations).Where("id = ?", id).Delete(&projects)
}

// Reset Password
func (projects *Projects) UpStatusProjects(status int) {
	sqlDB.Model(&projects).Update("new_status", status)
}

// Reset Password
func (projects *Projects) UpProjectsKey(key string) {
	sqlDB.Model(&projects).Update("key", key)
}

// Reset Password
func (projects *Projects) UpdateProjects(id string) {
	sqlDB.Model(&projects).
		Select("UserName", "Password", "AccNumber", "ColaAPI", "StatusJSON").
		Where("id = ?", id).
		Updates(&projects)
}
