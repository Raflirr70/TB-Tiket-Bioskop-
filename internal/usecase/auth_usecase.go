package usecase

import (
	"Project/internal/config"
	"Project/internal/domain/entity"
	"Project/internal/domain/repository"
	du "Project/internal/domain/usecase"
	"Project/pkg/helper"
	"errors"
	"time"
)

type AuthUsecaseImpl struct {
	roleRepository repository.RoleRepository
	userRepository repository.UserRepository
	config         *config.Config
}

func NewAuthUsecase(
	roleRepository repository.RoleRepository,
	userRepository repository.UserRepository,
	config *config.Config,
) du.AuthUsecase {
	return &AuthUsecaseImpl{
		roleRepository: roleRepository,
		userRepository: userRepository,
		config:         config,
	}
}

func (u *AuthUsecaseImpl) Login(req du.LoginRequest) (string, error) {
	user, err := u.userRepository.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("Invalid email or password")
	}

	if !helper.CheckPassword(req.Password, user.Password) {
		return "", errors.New("Invalid email or password")
	}

	token, err := helper.GenerateToken(user.ID, user.Email, user.Firstname, user.Lastname, user.RoleID, u.config.JWT)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *AuthUsecaseImpl) Register(req du.RegisterRequest) error {
	existingEmail, _ := u.userRepository.FindByEmail(req.Email)
	if existingEmail != nil {
		return errors.New("Email alredy registered")
	}

	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return err
	}
	user := &entity.User{
		RoleID:    1,
		Email:     req.Email,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Password:  hash,
		CreatedAt: time.Now(),
	}

	err = u.userRepository.Create(user)
	if err != nil {
		return err
	}
	return nil
}
