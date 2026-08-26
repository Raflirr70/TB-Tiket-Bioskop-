package entity

import "time"

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	FilmID    uint
	Ratting   int
	Value     string
	Time      time.Time
	CreatedAt time.Time
}
