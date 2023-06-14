package database

type IssuedNumber struct {
	ID           uint
	IssuedNumber float64
}

func (issued *IssuedNumber) UpCoinTotalNumber() {
	sqlDB.Save(&issued)
}
