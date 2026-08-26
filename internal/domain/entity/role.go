package entity

type Role struct {
	ID   uint `gorm:"primaryKey"`
	Name string

	Users []User
}
