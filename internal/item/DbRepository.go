package item

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

func (repo DBRepository) GetItemByOrderID(orderID uint) (Item, bool, error) {
	var item Item

	searchErr := repo.DBConnection.Where("order_id", orderID).Find(&item)

	if searchErr.Error != nil {
		return Item{}, false, searchErr.Error
	}

	if searchErr.RowsAffected == 0 {
		return Item{}, false, nil
	} else {
		return item, true, nil
	}
}

func (repo DBRepository) StoreItem(userID uint32, orderID uint, status string, accrual float64) (Item, error) {
	item := Item{
		OrderID: orderID,
		Bonus:   accrual,
		UserID:  userID,
		Status:  status,
	}

	if storeErr := repo.DBConnection.Create(&item).Error; storeErr != nil {
		return Item{}, storeErr
	}

	return item, nil
}

func (repo DBRepository) GetUserItems(userID uint32) ([]Item, error) {
	var items []Item

	err := repo.DBConnection.Where("user_id", userID).Find(&items).Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo DBRepository) GetUnprocessedItems() ([]Item, error) {
	var items []Item

	err := repo.DBConnection.
		Where("status IN ?", []string{StatusNew, StatusProcessing}).
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (repo DBRepository) UpdateItemStatusAndAccrualByOrderID(orderID uint, status string, accrual float64) error {
	return repo.DBConnection.Table("items").
		Where("order_id", orderID).
		Updates(map[string]interface{}{
			"status": status,
			"bonus":  accrual,
		}).Error
}
