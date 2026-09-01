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
		genres, err := u.genreRepo.GetByFilmID(film.ID)
		if err != nil {
			return nil, err
		}

		resultGenres := make([]du.GenreResponse, 0, len(genres))

		for _, genre := range genres {
			resultGenres = append(resultGenres, du.GenreResponse{
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
			Genres:    resultGenres,
		})
	}

	return result, nil
}
