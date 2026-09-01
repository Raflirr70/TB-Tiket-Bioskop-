package repository

import (
	"Project/internal/domain/entity"
	du "Project/internal/domain/usecase"
)

type GenreRepository interface {
	Create(film *entity.Genre) error
	Delete(id uint) error
	GetByFilmID(id uint) ([]*du.GenreResponse, error)
	Update(film *entity.Genre) error
}
