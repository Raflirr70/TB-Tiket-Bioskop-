package entity

import "time"

type Film struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Synopsis  string
	Duration  int
	Price     int
	Status    string
	IrlImg    string
	UpdatedAt time.Time
	CreatedAt time.Time

	Rattings  []Ratting
	Bookmarks []Bookmark
	Genres    []Genre `gorm:"many2many:genre_films"`
	Sources   []Source
	Schedules []Schedule
}
