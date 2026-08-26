package entity

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null"`
	Password  string
	Firstname string
	Lastname  string
	CreatedAt time.Time
}
