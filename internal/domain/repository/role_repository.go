package repository

import "Project/internal/domain/entity"

type RoleRepository interface {
	Create(role *entity.Role) error
	// FindByID(id uint) (*entity.Role, error)
	// FindByName(name string) (*entity.Role, error)
	// FindAll() ([]*entity.Role, error)
	// Update(role *entity.Role) error
	// Delete(id uint) error
}
