package usecase

import (
	"Project/internal/domain/entity"
	"Project/internal/domain/repository"
	du "Project/internal/domain/usecase"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseImp struct {
	userRepository repository.UserRepository
}

func NewUserUsecase(userRepository repository.UserRepository) du.UserUseCase {
	return &UserUsecaseImp{
		userRepository: userRepository,
	}
}

func (uc *UserUsecaseImp) Create(user *entity.User) error {
	if user.Email == "" {
		return errors.New("Masukan Email")
	}
	if user.Firstname == "" {
		return errors.New("Masukan Firstname")
	}
	if user.Lastname == "" {
		return errors.New("Masukan Lastname")
	}
	if user.Password == "" {
		return errors.New("Masukan Password")
	}

	has, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(has)
	return uc.userRepository.Create(user)
}

func (uc *UserUsecaseImp) Delete(id uint) error {
	return uc.userRepository.Delete(id)
}

func (uc *UserUsecaseImp) FindById(id uint) (*du.UserResponse, error) {
	return uc.userRepository.FindById(id)
}

func (uc *UserUsecaseImp) GetAll() ([]*du.UserResponse, error) {
	return uc.userRepository.GetAll()
}

func (uc *UserUsecaseImp) Update(user *entity.User, id uint) error {
	if id == 0 {
		return errors.New("ID user invalid")
	}

	user.ID = id

	return uc.userRepository.Update(user)
}
