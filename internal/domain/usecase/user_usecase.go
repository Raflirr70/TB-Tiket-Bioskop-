package usecase

import "time"

type UserUseCase interface {
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
