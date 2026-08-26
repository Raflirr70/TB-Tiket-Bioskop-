package entity

import "time"

type Transaction struct {
	ID            uint `gorm:"primaryKey"`
	Status        string
	TotalPrice    int
	PaymentMethod string
	Createdat     time.Time
}
