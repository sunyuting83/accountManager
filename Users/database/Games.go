package database

type Games struct {
	ID           uint `gorm:"primaryKey"`
	GameName     string
	BasePrice    float64
	UnitPrice    float64
	SingleNumber int64
}
