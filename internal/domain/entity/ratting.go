package entity

import "time"

type Ratting struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	FilmID    uint
	Value     int
	Time      time.Time
	CreatedAt time.Time
}
