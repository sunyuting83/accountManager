package database

type SplitManager struct {
	ID         uint `gorm:"primaryKey"`
	ManagerID  uint
	Proportion float64
}
