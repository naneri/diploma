package item

import (
	"gorm.io/gorm"
	"time"
)

type Item struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	OrderId   uint
	UserId    uint32
}
