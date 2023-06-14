package database

type PointsRecord struct {
	ID            uint `gorm:"primaryKey"`
	CoinManagerID uint
	CoinUsersID   uint
	Coin          float64
	CreatedAt     int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt     int64 `gorm:"autoUpdateTime:milli"`
}

func (record *PointsRecord) Insert() (err error) {
	sqlDB.Create(&record)
	return nil
}
