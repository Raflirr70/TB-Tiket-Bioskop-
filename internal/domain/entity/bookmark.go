package entity

import "time"

type Bookmark struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	FilmID    uint
	CreatedAt time.Time
}
