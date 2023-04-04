package database

type Games struct {
	ID        uint `gorm:"primaryKey"`
	GameName  string
	Projects  []Projects
	CreatedAt int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt int64 `gorm:"autoUpdateTime:milli"`
}

func (game *Games) Insert() (err error) {
	sqlDB.Create(&game)
	return nil
}

func CheckGamesName(gamename string) (game *Games, err error) {
	if err = sqlDB.First(&game, "game_name = ? ", gamename).Error; err != nil {
		return
	}
	return
}

// Delete Admin
func (game *Games) DeleteOne(id int64) {
	// time.Sleep(time.Duration(100) * time.Millisecond)
	sqlDB.Where("id = ?", id).Delete(&game)
}

// Check ID
func GameCheckID(id int64) (game *Games, err error) {
	if err = sqlDB.First(&game, "id = ?", id).Error; err != nil {
		return
	}
	return
}

// Get Count
func (game *Games) GetCount() (count int64, err error) {
	if err = sqlDB.Model(&game).Count(&count).Error; err != nil {
		return
	}
	return
}

// Admin List
func GetGamesList(page, Limit int) (game []*Games, err error) {
	p := makePage(page, Limit)
	if err = sqlDB.
		Order("id desc").
		Limit(Limit).Offset(p).
		Find(&game).Error; err != nil {
		return
	}
	return
}

func GetAllGamesList() (game []*Games, err error) {
	if err = sqlDB.
		Order("id desc").
		Find(&game).Error; err != nil {
		return
	}
	return
}
