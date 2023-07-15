package database

type Order struct {
	ID          uint   `gorm:"primaryKey"`
	OrderCode   string `gorm:"index"`
	NewStatus   int
	Coin        float64
	CoinUsersID uint
	Remarks     string
	Accounts    []Accounts
	CreatedAt   int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt   int64 `gorm:"autoUpdateTime:milli"`
}

func (order *Order) Insert() (id uint, err error) {
	result := sqlDB.Create(&order)
	if result.Error != nil {
		return 0, result.Error
	}
	return order.ID, nil
}

// Get Count
func (order *Order) GetOrdersCount(id uint) (count int64, err error) {
	if err = sqlDB.Model(&order).
		Where("coin_users_id = ?", id).
		Count(&count).Error; err != nil {
		return
	}
	return
}

func GetOrdersList(page, Limit int, id uint) (order *[]Order, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Where("coin_users_id = ?", id).
		Order("updated_at desc").
		Limit(Limit).Offset(p).
		Find(&order).Error; err != nil {
		return
	}
	return
}

func GetOrdersDetail(id string, userid uint) (order *Order, err error) {
	if err = sqlDB.
		Preload("Accounts").
		Preload("Accounts.Games").
		First(&order, "id = ? AND coin_users_id = ?", id, userid).Error; err != nil {
		return
	}
	return
}

func GetOrdersDetailForRefund(id string, uid uint) (order *Order, err error) {
	if err = sqlDB.
		First(&order, "id = ? AND new_status = ? AND coin_users_id = ?", id, 0, uid).Error; err != nil {
		return
	}
	return
}

func UpOrdersToRefunding(ID string, status int, Remarks string) (order []*Order, err error) {
	sqlDB.
		Model(&order).
		Where("id = ?", ID).
		UpdateColumns(Order{NewStatus: status, Remarks: Remarks})
	return
}
