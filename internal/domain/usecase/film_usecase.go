package usecase

import "time"

type FilmUseCase interface {
	GetAllFilm() ([]FilmResponse, error)
}

// DTO
type FilmResponse struct {
	ID        uint
	Name      string
	Synopsis  string
	Duration  int
	Price     int
	Status    string
	UpdatedAt time.Time
	CreatedAt time.Time

	Genres    []GenreResponse
	Schedules []ScheduleRespone
}

type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ScheduleRespone struct {
	ID        uint `gorm:"primaryKey"`
	FilmID    uint
	RoomID    uint
	Status    string
	Price     int
	Date      time.Time
	Time      time.Time
	CreatedAt time.Time
}
