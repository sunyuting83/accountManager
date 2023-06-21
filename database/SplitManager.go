package database

type SplitManager struct {
	ID        uint `gorm:"primaryKey"`
	ManagerID uint
	Percent   float64
}
