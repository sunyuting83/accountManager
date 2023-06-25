package database

type BlockChain struct {
	ID                 uint `gorm:"primaryKey"`
	UsersID            uint
	UsersCoin          float64
	UsersPercent       float64
	ManagerID          uint
	ManagerCoin        float64
	ManagerPercent     float64
	CoinManagerIDs     string
	CoinManagerCoin    float64
	CoinManagerPercent float64
	OrderID            uint
	CreatedAt          int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt          int64 `gorm:"autoUpdateTime:milli"`
}

func (blockchain *BlockChain) Insert() (err error) {
	sqlDB.Create(&blockchain)
	return nil
}
