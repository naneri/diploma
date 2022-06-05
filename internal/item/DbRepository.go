package item

import "gorm.io/gorm"

type DbRepository struct {
	DbConnection *gorm.DB
}

func InitDatabaseRepository(dbConnection *gorm.DB) *DbRepository {
	dbRepo := DbRepository{
		DbConnection: dbConnection,
	}

	return &dbRepo
}

func (repo DbRepository) GetItemByOrderId(orderId uint) (Item, bool, error) {
	var item Item

	searchErr := repo.DbConnection.Where("order_id", orderId).Find(&item)

	if searchErr.Error != nil {
		return Item{}, false, searchErr.Error
	}

	if searchErr.RowsAffected == 0 {
		return Item{}, false, nil
	} else {
		return item, true, nil
	}
}
