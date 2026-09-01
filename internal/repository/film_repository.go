package repository

import (
	"Project/internal/domain/entity"
	dr "Project/internal/domain/repository"
	du "Project/internal/domain/usecase"

	"gorm.io/gorm"
)

type FilmRepository struct {
	db *gorm.DB
}

func NewFilmRepository(db *gorm.DB) dr.FilmRepository {
	return &FilmRepository{
		db: db,
	}
}

func (r *FilmRepository) Create(film *entity.Film) error {
	return r.db.Create(film).Error
}

func (r *FilmRepository) Delete(id uint) error {
	return r.db.Delete(id).Error
}

func (r *FilmRepository) FindById(id uint) (*du.FilmResponse, error) {
	var film entity.Film
	err := r.db.First(&film, id).Error
	if err != nil {
		return nil, err
	}

	return &du.FilmResponse{
		ID:        film.ID,
		Name:      film.Name,
		Synopsis:  film.Synopsis,
		Duration:  film.Duration,
		Price:     film.Price,
		Status:    film.Status,
		UpdatedAt: film.UpdatedAt,
		CreatedAt: film.CreatedAt,
	}, nil
}

func (r *FilmRepository) GetAll() ([]*du.FilmResponse, error) {
	var films []entity.Film

	err := r.db.Find(&films).Error
	if err != nil {
		return nil, err
	}

	respon := make([]*du.FilmResponse, 0, len(films))

	for _, film := range films {
		respon = append(respon, &du.FilmResponse{
			ID:        film.ID,
			Name:      film.Name,
			Synopsis:  film.Synopsis,
			Duration:  film.Duration,
			Price:     film.Price,
			Status:    film.Status,
			UpdatedAt: film.UpdatedAt,
			CreatedAt: film.CreatedAt,
		})
	}

	return respon, nil
}

func (r *FilmRepository) Update(user *entity.Film) error {
	return r.db.Save(user).Error
}
