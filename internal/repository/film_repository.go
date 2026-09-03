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
		IrlImg:    film.IrlImg,
		UpdatedAt: film.UpdatedAt,
		CreatedAt: film.CreatedAt,
	}, nil
}

func (r *FilmRepository) GetAll(limit int, sort string) ([]entity.Film, error) {
	var films []entity.Film

	query := r.db.
		Preload("Genres").
		Preload("Schedules").
		Preload("Schedules.ScheduleSeats").
		Preload("Rattings")
	switch sort {
	case "ratting":
		query = query.Order("(SELECT COALESCE(AVG(value), 0) FROM rattings WHERE rattings.film_id = films.id) DESC")
	default:
		query = query.Order("films.created_at DESC")
	}
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&films).Error
	if err != nil {
		return nil, err
	}

	return films, nil
}

func (r *FilmRepository) Update(user *entity.Film) error {
	return r.db.Save(user).Error
}
