package repository

import (
	"Project/internal/domain/entity"
	dr "Project/internal/domain/repository"
	du "Project/internal/domain/usecase"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) dr.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(id).Error
}

func (r *UserRepository) FindById(id uint) (*du.UserResponse, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &du.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (r *UserRepository) GetAll() ([]*du.UserResponse, error) {
	var users []entity.User

	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	respon := make([]*du.UserResponse, 0, len(users))

	for _, user := range users {
		respon = append(respon, &du.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Password:  user.Password,
			Firstname: user.Firstname,
			Lastname:  user.Lastname,
			CreatedAt: user.CreatedAt,
		})
	}

	return respon, nil
}

func (r *UserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}
