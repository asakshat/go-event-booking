package models

import (
	"gorm.io/gorm"
)

type UserCart struct {
	gorm.Model
	EventID    uint    `gorm:"not null"`
	Event      Event   `gorm:"foreignKey:EventID"`
	Quantity   int     `gorm:"not null"`
	TotalPrice float64 `gorm:"type:decimal(10,2);not null"`
}
