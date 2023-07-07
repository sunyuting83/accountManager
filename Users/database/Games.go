package database

type Games struct {
	ID           uint `gorm:"primaryKey"`
	GameName     string
	BasePrice    float64
	UnitPrice    float64
	SingleNumber int64
}

// Account List
func GetGamesList() (games *[]Games, err error) {
	if err = sqlDB.
		Order("id desc").
		Find(&games).Error; err != nil {
		return
	}
	return
}
