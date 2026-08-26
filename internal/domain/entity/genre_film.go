package entity

type GenreFilm struct {
	ID      uint `gorm:"primaryKey"`
	FilmID  uint
	GenreID uint
}
