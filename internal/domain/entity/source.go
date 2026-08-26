package entity

type Source struct {
	ID          uint `gorm:"primaryKey"`
	FilmID      uint
	Source      string
	Description string
}
