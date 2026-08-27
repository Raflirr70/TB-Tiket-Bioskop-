package usecase

type AuthUsecase interface {
	Login(req LoginRequest) (token string, redireq string, err error)
	Register(req RegisterRequest) error
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Email           string `json:"email" form:"email" validate:"required,email"`
	Firstname       string `json:"firstname" form:"firstname" validate:"required"`
	Lastname        string `json:"lastname" form:"lastname" validate:"required"`
	Password        string `json:"password" form:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" validate:"required,eqfield=Password"`
}
