package services

import (
	"github.com/naneri/diploma/internal/item"
	"github.com/naneri/diploma/internal/user"
	"github.com/naneri/diploma/internal/withdrawal"
	"gorm.io/gorm"
	"log"
)

func RunMigrations(db *gorm.DB) {
	migrationErr := db.AutoMigrate(&user.User{}, &item.Item{}, &withdrawal.Withdrawal{})

	if migrationErr != nil {
		log.Fatalf("error migrating: %s", migrationErr.Error())
	}
}
