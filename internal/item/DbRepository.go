package item

import (
	"gorm.io/gorm"
)

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

func (repo DbRepository) StoreItem(userId uint32, orderId uint, status string, accrual float64) (Item, error) {
	item := Item{
		OrderId: orderId,
		Bonus:   accrual,
		UserId:  userId,
		Status:  status,
	}

	if storeErr := repo.DbConnection.Create(&item).Error; storeErr != nil {
		return Item{}, storeErr
	}

	return item, nil
}

func (repo DbRepository) GetUserItems(userId uint32) ([]Item, error) {
	var items []Item

	err := repo.DbConnection.Where("user_id", userId).Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo DbRepository) GetUnprocessedItems() ([]Item, error) {
	var items []Item

	err := repo.DbConnection.
		Where("status IN ?", []string{StatusNew, StatusProcessing}).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo DbRepository) UpdateItemStatusAndAccrualByOrderId(orderId uint, status string, accrual float64) error {
	return repo.DbConnection.Table("items").
		Where("order_id", orderId).
		Updates(map[string]interface{}{
			"status": status,
			"bonus":  accrual,
		}).Error
}
