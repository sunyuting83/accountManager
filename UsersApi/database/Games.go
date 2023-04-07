package database

type Games struct {
	ID        uint `gorm:"primaryKey"`
	GameName  string
	Projects  []Projects
	CreatedAt int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}
