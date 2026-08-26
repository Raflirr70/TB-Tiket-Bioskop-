package repository

import (
	"Project/internal/domain/entity"
	du "Project/internal/domain/usecase"
)

type UserRepository interface {
	Create(user *entity.User) error
	Delete(id uint) error
	FindById(id uint) (*du.UserResponse, error)
	GetAll() ([]*du.UserResponse, error)
	Update(user *entity.User) error
}
