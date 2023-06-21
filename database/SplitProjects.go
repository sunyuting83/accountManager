package database

type SplitProjects struct {
	ID      uint `gorm:"primaryKey"`
	Percent float64
}
