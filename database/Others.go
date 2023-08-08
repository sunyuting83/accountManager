package database

type Others struct {
	ID         uint `gorm:"primaryKey"`
	ProjectsID uint
	Account    string
	Password   string
	JsonData   map[string]interface{} `gorm:"column:json_data;type:jsonb"`
	NewStatus  int                    `gorm:"index"`
	Exptime    int64
	Remarks    string
	CreatedAt  int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt  int64 `gorm:"autoUpdateTime:milli"`
}

// Account List
func GetOthersList(ProjectsID string) (others []Others, err error) {
	if err = sqlDB.
		Where("projects_id = ?", ProjectsID).
		Order("updated_at desc").
		Find(&others).Error; err != nil {
		return
	}
	return
}

func OthersBatches(others []Others) {
	sqlDB.Create(&others)
}
func OthersInBatches(others []Others) {
	sqlDB.CreateInBatches(&others, 1000)
}
