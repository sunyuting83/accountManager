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
	OrderID       uint
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
	Projects      Projects
	Games         Games `gorm:"foreignKey:GameID"`
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

func WithGameID(GameID uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if GameID > 0 {
			db.Where("game_id = ?", GameID)
		}
		return db.Where("")
	}
}

// Get Count
func (accounts *Accounts) GetCountWithSellStatus(gameid uint) (count int64, err error) {
	if err = sqlDB.Model(&accounts).
		Where("sell_status = 1").
		Scopes(WithGameID(gameid)).
		Count(&count).Error; err != nil {
		return
	}
	return
}

// Account List
func GetAccountList(page, Limit int, GameID uint) (accounts *[]Accounts, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("sell_status = 1 AND new_status != 108").
		Scopes(WithGameID(GameID)).
		Preload("Games").
		Order("today_gold desc").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}

// Account List
func GetAccountsWithIn(IDs []int) (accounts *[]Accounts, err error) {
	if err = sqlDB.
		Where("sell_status = 1 AND new_status != 108 AND id IN ?", IDs).
		Preload("Games").
		Preload("Projects.Users").
		Find(&accounts).Error; err != nil {
		return
	}
	return
}
func UpAccountsWithIn(IDs []string) (accounts []*Accounts, err error) {
	sqlDB.
		Model(&accounts).
		Clauses(clause.Returning{}).
		Where("id IN ?", IDs).
		UpdateColumns(Accounts{SellStatus: 2, NewStatus: 108})
	return
}

func UpAccountsToRefund(ID []int) (accounts []*Accounts, err error) {
	sqlDB.
		Model(&accounts).
		Where("id = IN", ID).
		UpdateColumns(Accounts{SellStatus: 120})
	return
}

func UpOrderIDForAccountsWithIn(IDs []string, OrderID uint) {
	sqlDB.
		Model(&Accounts{}).
		Where("id IN ?", IDs).
		UpdateColumns(Accounts{OrderID: OrderID})
}

// Account List
func GetFailedAccountsWithIn(IDs []int) (accounts *[]Accounts, err error) {
	if err = sqlDB.
		Where("id IN ?", IDs).
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

func (accounts *Accounts) GetCountUseScopesB(filter *utils.Filter, page, Limit int, GameID uint) (count int64, err error) {
	if err = sqlDB.Model(&accounts).
		Where("sell_status = 1 AND new_status != 108").
		Scopes(WithGameID(GameID)).
		Preload("Games").
		Scopes(MinGold(filter.MinGold)).
		Scopes(MaxGold(filter.MaxGold)).
		Scopes(Multiple(filter.Multiple)).
		Scopes(Diamond(filter.Diamond)).
		Scopes(Crazy(filter.Crazy)).
		Scopes(Cold(filter.Cold)).
		Scopes(Precise(filter.Precise)).
		Count(&count).Error; err != nil {
		return
	}
	return
}

func GetDataUseScopesB(filter *utils.Filter, page, Limit int, GameID uint) (accounts *[]Accounts, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("sell_status = 1 AND new_status != 108").
		Scopes(WithGameID(GameID)).
		Preload("Games").
		Scopes(MinGold(filter.MinGold)).
		Scopes(MaxGold(filter.MaxGold)).
		Scopes(Multiple(filter.Multiple)).
		Scopes(Diamond(filter.Diamond)).
		Scopes(Crazy(filter.Crazy)).
		Scopes(Cold(filter.Cold)).
		Scopes(Precise(filter.Precise)).
		Order("today_gold DESC").
		Limit(Limit).Offset(p).
		Find(&accounts).Error; err != nil {
		return
	}
	return
}
