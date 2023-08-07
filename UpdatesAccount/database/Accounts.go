package database

type Accounts struct {
	ID            uint `gorm:"primaryKey"`
	ProjectsID    uint
	GameID        uint
	ComputID      uint
	OrderID       *uint
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
	CreatedAt     int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt     int64 `gorm:"autoUpdateTime:milli"`
}

func UpdateAccountGamesID() (err error) {
	var account []Accounts
	if err := sqlDB.
		Model(&Accounts{}).
		Where("new_status != 8 AND new_status != 9 AND new_status != 10").
		Preload("Projects.Games").
		Find(&account).Error; err != nil {
		return err
	}

	for i := range account {
		account[i].GameID = account[i].Projects.GamesID
	}

	if err := sqlDB.Save(&account).Error; err != nil {
		return err
	}

	return nil
}

func GetCount() (count int64, err error) {
	if err := sqlDB.
		Model(&Accounts{}).
		Where("new_status IN (8,9,10)").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return
}
