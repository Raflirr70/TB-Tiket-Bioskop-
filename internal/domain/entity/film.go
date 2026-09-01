package entity

import "time"

type Film struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Synopsis  string
	Duration  int
	Price     int
	Status    string
	UpdatedAt time.Time
	CreatedAt time.Time

	Comments  []Comment
	Bookmarks []Bookmark
	Genres    []Genre `gorm:"many2many:genre_films;"`
	Sources   []Source
	Schedules []Schedule
}
