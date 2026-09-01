package repository

import (
	"Project/internal/domain/entity"
	du "Project/internal/domain/usecase"
)

type FilmRepository interface {
	Create(film *entity.Film) error
	Delete(id uint) error
	FindById(id uint) (*du.FilmResponse, error)
	GetAll() ([]du.FilmResponse, error)
	Update(film *entity.Film) error
}
