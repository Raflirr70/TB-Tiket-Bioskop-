package usecase

import (
	"Project/internal/domain/repository"
	d "Project/internal/domain/usecase"
)

type UserUsecaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) d.UserUseCase {
	return &UserUsecaseImpl{
		userRepository: userRepository,
	}
}
