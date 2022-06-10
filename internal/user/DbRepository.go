package user

import (
	"gorm.io/gorm"
	"sync"
)

type DBRepository struct {
	// not sure if this is correct to add a mutex.Lock() features to this repo, as in a high-load system, this will give a huge overhead. I would prefer to only lock access to a single user balance.
	Access       sync.Mutex
	DBConnection *gorm.DB
}

func InitDatabaseRepository(dbConnection *gorm.DB) *DBRepository {
	dbRepo := DBRepository{
		DBConnection: dbConnection,
	}

	return &dbRepo
}

func (dbRepo *DBRepository) Save(login, hashedPass string) (User, error) {
	dbRepo.Access.Lock()
	defer dbRepo.Access.Unlock()
	var user User

	searchErr := dbRepo.DBConnection.Where("login", login).Find(&user).Error
	if searchErr != nil {
		return User{}, searchErr
	}

	if user.Login == login {
		return User{}, &AlreadyExistsError{login}
	}

	user.Login = login
	user.Password = hashedPass
	saveErr := dbRepo.DBConnection.Create(&user).Error

	if saveErr != nil {
		return User{}, saveErr
	}

	return user, nil
}

func (dbRepo *DBRepository) Find(login string) (User, bool, error) {
	dbRepo.Access.Lock()
	defer dbRepo.Access.Unlock()
	var user User

	searchErr := dbRepo.DBConnection.Where("login", login).Find(&user)

	if searchErr.Error != nil {
		return User{}, false, searchErr.Error
	}

	if searchErr.RowsAffected == 0 {
		return User{}, false, nil
	} else {
		return user, true, nil
	}
}

func (dbRepo *DBRepository) FindUserByID(id uint32) (User, error) {
	dbRepo.Access.Lock()
	defer dbRepo.Access.Unlock()
	var user User

	searchErr := dbRepo.DBConnection.Where("id", id).First(&user).Error

	if searchErr != nil {
		return User{}, searchErr
	}

	return user, nil
}

func (dbRepo *DBRepository) UpdateUserBalance(userID uint32, bonus float64) error {
	dbRepo.Access.Lock()
	defer dbRepo.Access.Unlock()
	var user User

	if searchErr := dbRepo.DBConnection.Where("id", userID).First(&user).Error; searchErr != nil {
		return searchErr
	}

	user.Balance += bonus

	if saveErr := dbRepo.DBConnection.Save(&user).Error; saveErr != nil {
		return saveErr
	}

	return nil
}
