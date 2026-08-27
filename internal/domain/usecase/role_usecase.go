package usecase

type RoleUsecase interface {
	GetRoleByID(id uint) (*RoleResponse, error)
	CreateRole(req CreateRoleRequest) (*RoleResponse, error)
}

type CreateRoleRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
