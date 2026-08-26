package usecase

import (
	"Project/internal/domain/repository"
	du "Project/internal/domain/usecase"
)

type UserUsecaseImp struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) du.UserUseCase {
	return &UserUsecaseImp{
		userRepository: userRepository,
	}
}
