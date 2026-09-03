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

func (u *FilmUsecaseImpl) GetAllFilm(limit int, sort string) ([]du.FilmResponse, error) {
	films, err := u.filmRepo.GetAll(limit, sort)
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
		var totalRating int

		for _, rating := range film.Rattings {
			totalRating += rating.Value
		}

		var averageRating float32

		if len(film.Rattings) > 0 {
			averageRating = float32(totalRating) / float32(len(film.Rattings))
		}

		schedules := make([]du.ScheduleResponse, 0, len(film.Schedules))
		for _, schedule := range film.Schedules {
			seats := make([]du.ScheduleSeatResponse, 0, len(schedule.ScheduleSeats))

			for _, seat := range schedule.ScheduleSeats {
				seats = append(seats, du.ScheduleSeatResponse{
					ID:         seat.ID,
					ScheduleID: seat.ScheduleID,
					SeatID:     seat.SeatID,
					Status:     seat.Status,
					Time:       seat.Time,
				})
			}
			schedules = append(schedules, du.ScheduleResponse{
				ID:             schedule.ID,
				FilmID:         schedule.FilmID,
				RoomID:         schedule.RoomID,
				Status:         schedule.Status,
				Date:           schedule.Date,
				Time:           schedule.Time,
				CreatedAt:      schedule.CreatedAt,
				SchedulesSeats: seats,
			})
		}

		result = append(result, du.FilmResponse{
			ID:        film.ID,
			Name:      film.Name,
			Synopsis:  film.Synopsis,
			Duration:  film.Duration,
			Price:     film.Price,
			Status:    film.Status,
			IrlImg:    film.IrlImg,
			UpdatedAt: film.UpdatedAt,
			CreatedAt: film.CreatedAt,
			Genres:    genres,
			Schedules: schedules,
			Ratting:   averageRating,
		})
	}

	return result, nil
}
