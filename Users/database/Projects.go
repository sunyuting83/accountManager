package database

type Projects struct {
	ID      uint `gorm:"primaryKey"`
	UsersID uint
	Users   Users
}
