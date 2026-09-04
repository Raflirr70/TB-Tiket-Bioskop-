package repository

import "Project/internal/domain/entity"

type RoomRepository interface {
	GetAll() ([]entity.Room, error)
	Create(room *entity.Room) error
	FindById(id uint) (*entity.Room, error)
	Update(room *entity.Room) error
	SoftDelete(id uint) error
	ReplaceSeats(roomID uint, seats []entity.Seat) error
	SyncSeats(roomID uint, requestedSeatNames []string) error
	HasActiveSchedules(roomID uint) (bool, error)
	CountActiveSchedules(roomID uint) (int64, error)
}
