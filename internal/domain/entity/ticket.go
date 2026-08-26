package entity

import "time"

type Ticket struct {
	ID             uint `gorm:"primaryKey"`
	UserID         uint
	ScheduleSeatID uint
	TransactionID  uint
	Time           time.Time
	CreatedAt      time.Time
}
