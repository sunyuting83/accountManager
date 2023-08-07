package database

type Games struct {
	ID           uint `gorm:"primaryKey"`
	GameName     string
	Projects     []Projects
	Count        int64
	AliveCount   int64
	BasePrice    float64
	UnitPrice    float64
	SingleNumber int64
	CreatedAt    int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt    int64 `gorm:"autoUpdateTime:milli"`
}
