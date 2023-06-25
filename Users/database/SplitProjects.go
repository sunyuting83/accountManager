package database

type SplitProjects struct {
	ID      uint `gorm:"primaryKey"`
	Percent float64
	Manager string
}

func GetSplittedPercent() (splitProjects *SplitProjects, err error) {
	if err = sqlDB.First(&splitProjects, "id = ? ", 1).Error; err != nil {
		return
	}
	return
}
