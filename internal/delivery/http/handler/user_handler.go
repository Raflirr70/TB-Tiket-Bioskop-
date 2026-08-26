package handler

import (
	"Project/internal/domain/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserUsecase usecase.UserUseCase
}

func NewUserHandler(userUsecase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.UserUsecase.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": err.Error(),
		})
		return
	}
	c.HTML(http.StatusOK, "users.html", gin.H{
		"title": "Users",
		"users": users,
	})
}
