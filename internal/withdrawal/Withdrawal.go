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
	UserId    uint32
	Sum       float64
	OrderId   uint
}
