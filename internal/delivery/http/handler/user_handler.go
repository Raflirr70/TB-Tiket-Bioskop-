package handler

import (
	"Project/internal/domain/usecase"
)

type UserHandler struct {
	UserUsecase usecase.UserUseCase
}

func NewUserHandler(userUsecase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}
