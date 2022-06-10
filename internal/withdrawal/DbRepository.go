package withdrawal

import (
	"gorm.io/gorm"
)

type DBRepository struct {
	DBConnection *gorm.DB
}

func InitDatabaseRepository(dbConnection *gorm.DB) *DBRepository {
	dbRepo := DBRepository{
		DBConnection: dbConnection,
	}

	return &dbRepo
}

func (repo DBRepository) ListUserWithdrawals(UserID uint32) ([]Withdrawal, error) {
	var withdrawals []Withdrawal

	err := repo.DBConnection.Where("user_id", UserID).Find(&withdrawals).Error

	return withdrawals, err
}
