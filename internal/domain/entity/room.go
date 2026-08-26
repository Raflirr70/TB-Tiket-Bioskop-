package entity

type Room struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Capasity int

	Schedules []Schedule
	Seats     []Seat
}
