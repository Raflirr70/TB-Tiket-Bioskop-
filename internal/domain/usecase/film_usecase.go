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

	Genres []GenreResponse
}

type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
