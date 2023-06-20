package database

type Users struct {
	ID        uint `gorm:"primaryKey"`
	ManagerID uint
}
