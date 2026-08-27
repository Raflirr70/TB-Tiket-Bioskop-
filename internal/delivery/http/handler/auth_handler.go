package handler

import (
	du "Project/internal/domain/usecase"
	"Project/pkg/response"
	"Project/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase du.AuthUsecase
}

func NewAuthHandler(authUsecase du.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	response.Success(c, http.StatusOK, "Logged out successfully")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req du.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.authUsecase.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Set cookie "token" for browser access
	c.SetCookie("token", token, 3600*24, "/", "", false, true)

	response.Success(c, http.StatusOK, gin.H{
		"token":    token,
		"redirect": "/dashboard",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req du.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := validator.Validate(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.authUsecase.Register(req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}
