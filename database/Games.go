package database

import "gorm.io/gorm"

type Games struct {
	ID         uint `gorm:"primaryKey"`
	GameName   string
	Projects   []Projects
	Count      int64
	AliveCount int64
	CreatedAt  int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt  int64 `gorm:"autoUpdateTime:milli"`
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

func GetAllGames() (game []*Games, err error) {
	if err = sqlDB.
		Preload("Projects", func(db *gorm.DB) *gorm.DB {
			return sqlDB.Select("ID", "GamesID", "StatusJSON", "UsersID")
		}).
		Select("ID", "GameName").
		Order("id desc").
		Find(&game).Error; err != nil {
		return
	}
	return
}

func GetGame(id int64) (game *Games, err error) {
	if err = sqlDB.
		Preload("Projects").
		Preload("Projects.Users").
		Select("ID", "GameName").
		Where("id", id).
		Order("id desc").
		Find(&game).Error; err != nil {
		return
	}
	return
}

/*
, func(db *gorm.DB) *gorm.DB {
	return sqlDB.Select("ID", "GamesID", "StatusJSON", "UsersID")
}
*/
