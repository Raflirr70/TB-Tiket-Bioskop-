package entity

import "gorm.io/gorm"

type Seat struct {
	ID        uint `gorm:"primaryKey"`
	RoomID    uint
	Name      string
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
