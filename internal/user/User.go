package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint32 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Login     string
	Password  string
	WithDrawn float64 `sql:"DEFAULT:0"`
	Balance   float64 `sql:"DEFAULT:0"`
}
