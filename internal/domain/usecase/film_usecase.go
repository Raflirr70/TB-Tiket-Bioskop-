package usecase

import "time"

type FilmUseCase interface {
	GetAllFilm(limit int, sort string) ([]FilmResponse, error)
}

// DTO
type FilmResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Synopsis  string    `json:"synopsis"`
	Duration  int       `json:"duration"`
	Price     int       `json:"price"`
	Status    string    `json:"status"`
	IrlImg    string    `json:"irl_img"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Ratting   float32   `json:"ratting"`

	Genres    []GenreResponse    `json:"genres"`
	Schedules []ScheduleResponse `json:"schedules"`
}

type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ScheduleResponse struct {
	ID        uint      `json:"id"`
	FilmID    uint      `json:"film_id"`
	RoomID    uint      `json:"room_id"`
	Status    string    `json:"status"`
	Price     int       `json:"price"`
	Date      time.Time `json:"date"`
	Time      time.Time `json:"time"`
	CreatedAt time.Time `json:"created_at"`

	SchedulesSeats []ScheduleSeatResponse `json:"schedules_seats"`
}

type ScheduleSeatResponse struct {
	ID         uint      `json:"id"`
	ScheduleID uint      `json:"schedule_id"`
	SeatID     uint      `json:"seat_id"`
	Status     string    `json:"status"`
	Time       time.Time `json:"time"`
}

type RattingResponse struct {
	ID    uint `json:"id"`
	Value int  `json:"value"`
}
