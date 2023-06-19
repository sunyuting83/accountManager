package database

import (
	"colaAPI/Users/utils"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Accounts struct {
	ID            uint `gorm:"primaryKey"`
	ProjectsID    uint
	GameID        uint
	ComputID      uint
	PhoneNumber   string
	PhonePassword string
	UserName      string
	Password      string
	Cover         string
	NewStatus     int `gorm:"index"`
	SellStatus    int `gorm:"index;default:0"`
	TodayGold     int64
	YesterdayGold int64
	Multiple      int64
	Diamond       int
	Crazy         int
	Precise       int
	Cold          int
	Exptime       int64
	Price         float64
	Remarks       string
	CreatedAt     int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt     int64 `gorm:"autoUpdateTime:milli"`
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

// Get Count
func (accounts *Accounts) GetCount(ProjectsID, Status string) (count int64, err error) {
	if err = sqlDB.Model(&accounts).Where("projects_id = ? and new_status = ?", ProjectsID, Status).Count(&count).Error; err != nil {
		return
	}
	return
}

// Account List
func (account *Accounts) GetInList(ProjectsID string, statusList []string, page, Limit int) (accounts []*Accounts, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("projects_id = ? and new_status IN ?", ProjectsID, statusList).
		Order("today_gold DESC").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

// Account List
func GetAccountList(page, Limit int, ProjectsID, Status string) (accounts *[]Accounts, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("projects_id = ? and new_status = ?", ProjectsID, Status).
		Order("updated_at desc").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

func MinGold(MinGold int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if MinGold > 0 {
			return db.Where("today_gold >= ?", MinGold)
		}
		return db.Where("")
	}
}
func MaxGold(MaxGold int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if MaxGold > 0 {
			return db.Where("today_gold <= ?", MaxGold)
		}
		return db.Where("")
	}
}
func Multiple(Multiple int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Multiple > 0 {
			return db.Where("multiple >= ?", Multiple)
		}
		return db.Where("")
	}
}
func Diamond(Diamond int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Diamond > 0 {
			return db.Where("diamond >= ?", Diamond)
		}
		return db.Where("")
	}
}
func Crazy(Crazy int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Crazy > 0 {
			return db.Where("crazy >= ?", Crazy)
		}
		return db.Where("")
	}
}
func Cold(Cold int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Cold > 0 {
			return db.Where("cold >= ?", Cold)
		}
		return db.Where("")
	}
}
func Precise(Precise int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if Precise > 0 {
			return db.Where("precise >= ?", Precise)
		}
		return db.Where("")
	}
}
func HasStatus(hasStatus []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("new_status IN (?)", hasStatus)
	}
}

func SetSellUseIn(IDs []int) (accounts []*Accounts, err error) {
	sqlDB.
		Model(&accounts).
		Clauses(clause.Returning{}).
		Where("id IN ?", IDs).
		Update("sell_status", "1")
	return
}

func GetDataUseScopes1(filter utils.Filter, hasStatus []string, projectsID string) (accounts []Accounts, err error) {
	if err = sqlDB.
		Where("projects_id = ?", projectsID).
		Scopes(HasStatus(hasStatus)).
		Scopes(MinGold(filter.MinGold)).
		Scopes(MaxGold(filter.MaxGold)).
		Scopes(Multiple(filter.Multiple)).
		Scopes(Diamond(filter.Diamond)).
		Scopes(Crazy(filter.Crazy)).
		Scopes(Cold(filter.Cold)).
		Scopes(Precise(filter.Precise)).
		Order("today_gold DESC").
		Find(&accounts).Error; err != nil {
		return
	}
	return
}
