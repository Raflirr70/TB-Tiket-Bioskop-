package entity

import "time"

type Schedule struct {
	ID        uint `gorm:"primaryKey"`
	FilmID    uint
	RoomID    uint
	Status    string
	Date      time.Time
	Time      time.Time
	CreatedAt time.Time

	ScheduleSeats []ScheduleSeat
}
