package entity

import "gorm.io/gorm"

type Room struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Capasity  int
	Status    string         `gorm:"default:ready"` // ready | not_ready | deleted
	DeletedAt gorm.DeletedAt // soft delete

	Schedules []Schedule
	Seats     []Seat
}
