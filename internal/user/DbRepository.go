package user

import (
	"gorm.io/gorm"
	"sync"
)

type DbRepository struct {
	// not sure if this is correct to add a mutex.Lock() features to this repo, as in a high-load system, this will give a huge overhead. I would prefer to only lock access to a single user balance.
	Access       sync.Mutex
	DbConnection *gorm.DB
}

func InitDatabaseRepository(dbConnection *gorm.DB) *DbRepository {
	dbRepo := DbRepository{
		DbConnection: dbConnection,
	}

	return &dbRepo
}

func (dbRepo *DbRepository) Save(login, hashedPass string) (User, error) {
	dbRepo.Access.Lock()
	defer dbRepo.Access.Unlock()
	var user User

	searchErr := dbRepo.DbConnection.Where("login", login).Find(&user).Error
	if searchErr != nil {
		return User{}, searchErr
	}

	if user.Login == login {
		return User{}, &AlreadyExistsError{login}
	}

	user.Login = login
	user.Password = hashedPass
	saveErr := dbRepo.DbConnection.Create(&user).Error

	if saveErr != nil {
		return User{}, saveErr
	}

	return user, nil
}

func (dbRepo *DbRepository) Find(login string) (User, bool, error) {
	dbRepo.Access.Lock()
	defer dbRepo.Access.Unlock()
	var user User

	searchErr := dbRepo.DbConnection.Where("login", login).Find(&user)

	if searchErr.Error != nil {
		return User{}, false, searchErr.Error
	}

	if searchErr.RowsAffected == 0 {
		return User{}, false, nil
	} else {
		return user, true, nil
	}
}

func (dbRepo *DbRepository) FindUserById(id uint32) (User, error) {
	dbRepo.Access.Lock()
	defer dbRepo.Access.Unlock()
	var user User

	searchErr := dbRepo.DbConnection.Where("id", id).First(&user).Error

	if searchErr != nil {
		return User{}, searchErr
	}

	return user, nil
}

func (dbRepo *DbRepository) UpdateUserBalance(userId uint32, bonus float64) error {
	dbRepo.Access.Lock()
	defer dbRepo.Access.Unlock()
	var user User

	if searchErr := dbRepo.DbConnection.Where("id", userId).First(&user).Error; searchErr != nil {
		return searchErr
	}

	user.Balance += bonus

	if saveErr := dbRepo.DbConnection.Save(&user).Error; saveErr != nil {
		return saveErr
	}

	return nil
}
