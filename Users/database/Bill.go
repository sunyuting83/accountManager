package database

type Bill struct {
	ID              uint `gorm:"primaryKey"`
	CoinUsersID     *uint
	FormCoinUsersID *uint `gorm:"foreignKey:CoinUsersID"`
	OrderID         *uint
	Coin            float64
	NewStatus       int
	FormCoinUsers   CoinUsers
	Months          string `gorm:"index"`
	CreatedAt       int64  `gorm:"autoUpdateTime:milli"`
	UpdatedAt       int64  `gorm:"autoUpdateTime:milli"`
}

func (bill *Bill) Insert() (err error) {
	sqlDB.Create(&bill)
	return nil
}

// Get Count
func (bill *Bill) GetLedgerCount(id uint) (count int64, err error) {
	if err = sqlDB.Model(&bill).
		Where("coin_users_id = ?", id).
		Count(&count).Error; err != nil {
		return
	}
	return
}

func GetLedger(page, Limit int, userid uint) (bills *[]Bill, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("coin_users_id = ?", userid).
		Preload("FormCoinUsers").
		Order("created_at DESC").
		Limit(Limit).Offset(p).
		Find(&bills).Error; err != nil {
		return
	}
	return
}
