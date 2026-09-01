package repository

import (
	"Project/internal/domain/entity"
	dr "Project/internal/domain/repository"
	du "Project/internal/domain/usecase"

	"gorm.io/gorm"
)

type GenreRepository struct {
	db *gorm.DB
}

func NewGenreRepository(db *gorm.DB) dr.GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

func (r *GenreRepository) Create(genre *entity.Genre) error {
	return r.db.Create(genre).Error
}

func (r *GenreRepository) Delete(id uint) error {
	return r.db.Delete(id).Error
}

func (r *GenreRepository) GetByFilmID(id uint) ([]*du.GenreResponse, error) {
	var films []entity.Genre

	err := r.db.Find(&films).Error
	if err != nil {
		return nil, err
	}

	respon := make([]*du.GenreResponse, 0, len(films))
	for _, film := range films {
		respon = append(respon, &du.GenreResponse{
			ID:   film.ID,
			Name: film.Name,
		})
	}

	return respon, nil
}

func (r *GenreRepository) Update(genre *entity.Genre) error {
	return r.db.Save(genre).Error
}
