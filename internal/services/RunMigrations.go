package services

import (
	"github.com/naneri/diploma/internal/item"
	"github.com/naneri/diploma/internal/user"
	"gorm.io/gorm"
	"log"
)

func RunMigrations(db *gorm.DB) {
	migrationErr := db.AutoMigrate(&user.User{}, &item.Item{})

	if migrationErr != nil {
		log.Fatalf("error migrating: %s", migrationErr.Error())
	}
}
