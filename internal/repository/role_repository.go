package repository

import (
	"Project/internal/domain/entity"
	dr "Project/internal/domain/repository"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) dr.RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}
