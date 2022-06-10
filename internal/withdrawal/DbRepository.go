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
