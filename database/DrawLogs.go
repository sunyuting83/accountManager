package database

type DrawLogs struct {
	ID         uint `gorm:"primaryKey"`
	ProjectsID uint
	DrawUser   string
	LogName    string
	Data       string
	CreatedAt  int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt  int64 `gorm:"autoUpdateTime:milli"`
}

func (draw *DrawLogs) AddDrawLogs() {
	sqlDB.Create(&draw)
}
