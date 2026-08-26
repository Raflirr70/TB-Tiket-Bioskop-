package usecase

import (
	"Project/internal/domain/entity"
	"time"
)

type UserUseCase interface {
	Create(user *entity.User) error
	Delete(id uint) error
	FindById(id uint) (*UserResponse, error)
	GetAll() ([]*UserResponse, error)
	Update(user *entity.User, id uint) error
}

// DTO
type UserResponse struct {
	ID        uint
	Email     string
	Password  string
	Firstname string
	Lastname  string
	CreatedAt time.Time
}
