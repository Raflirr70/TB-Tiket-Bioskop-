package entity

import "time"

type Film struct {
	ID          uint `gorm:"primaryKey"`
	Name        string
	Description string
	Duration    string
	Status      string
	UpdatedAt   time.Time
	CreatedAt   time.Time

	Comments   []Comment
	Bookmarks  []Bookmark
	GenreFilms []GenreFilm
	Sources    []Source
	Schedules  []Schedule
}
