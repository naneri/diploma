package withdrawal

import (
	"gorm.io/gorm"
	"time"
)

type Withdrawal struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserID    uint32
	Sum       float64
	OrderID   uint
}
