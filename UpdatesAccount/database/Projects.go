package database

type Projects struct {
	ID           uint `gorm:"primaryKey"`
	UsersID      uint
	GamesID      uint
	ProjectsName string
	UserName     string
	Password     string
	StatusJSON   string
	NewStatus    int `gorm:"index"`
	Accounts     []Accounts
	Key          string
	Remarks      string
	AccNumber    int
	ColaAPI      bool
	Games        Games
	CreatedAt    int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt    int64 `gorm:"autoUpdateTime:milli"`
}
