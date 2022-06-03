package user

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

func (dbRepo *DbRepository) Save(login, hashedPass string) (User, error) {
	var user User

	searchErr := dbRepo.DbConnection.Where("login", login).Find(&user).Error
	if searchErr != nil {
		return User{}, searchErr
	}

	if user.Username == login {
		return User{}, &alreadyExists{login}
	}

	user.Username = login
	user.Password = hashedPass
	saveErr := dbRepo.DbConnection.Create(&user).Error

	if saveErr != nil {
		return User{}, saveErr
	}

	return user, nil
}
