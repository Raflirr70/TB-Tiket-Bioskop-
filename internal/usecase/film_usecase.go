package usecase

import (
	"Project/internal/config"
	"Project/internal/domain/repository"
	du "Project/internal/domain/usecase"
)

type FilmUsecaseImpl struct {
	filmRepo  repository.FilmRepository
	genreRepo repository.GenreRepository
	cg        *config.Config
}

func NewFilmUsecase(
	filmRepo repository.FilmRepository,
	genreRepo repository.GenreRepository,
	cg *config.Config,
) du.FilmUseCase {
	return &FilmUsecaseImpl{
		filmRepo:  filmRepo,
		genreRepo: genreRepo,
		cg:        cg,
	}
}

func (u *FilmUsecaseImpl) GetAllFilm() ([]du.FilmResponse, error) {
	films, err := u.filmRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]du.FilmResponse, 0, len(films))

	for _, film := range films {
		genres := make([]du.GenreResponse, 0, len(film.Genres))

		for _, genre := range film.Genres {
			genres = append(genres, du.GenreResponse{
				ID:   genre.ID,
				Name: genre.Name,
			})
		}

		result = append(result, du.FilmResponse{
			ID:        film.ID,
			Name:      film.Name,
			Synopsis:  film.Synopsis,
			Duration:  film.Duration,
			Price:     film.Price,
			Status:    film.Status,
			UpdatedAt: film.UpdatedAt,
			CreatedAt: film.CreatedAt,
			Genres:    genres,
		})
	}

	return result, nil
}
