package entity

type Genre struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
