package entity

type Seat struct {
	ID     uint `gorm:"primaryKey"`
	RoomID uint
	Name   string
}
