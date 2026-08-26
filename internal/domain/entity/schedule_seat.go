package entity

import "time"

type ScheduleSeat struct {
	ID         uint `gorm:"primaryKey"`
	ScheduleID uint
	SeatID     uint
	Status     string
	Time       time.Time
	Update     time.Time

	Ticket Ticket
}
