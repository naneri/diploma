package item

import (
	"gorm.io/gorm"
	"time"
)

const StatusNew = "NEW"
const StatusProcessing = "PROCESSING"
const StatusInavalid = "INVALID"
const StatusProcessed = "PROCESSED"

type Item struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	OrderID   uint
	Bonus     float64
	UserID    uint32
	Status    string
}
