package handler

import (
	"Project/internal/config"
	du "Project/internal/domain/usecase"
	"Project/pkg/response"
	"Project/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase du.AuthUsecase
	cg          *config.Config
}

func NewAuthHandler(authUsecase du.AuthUsecase, cg *config.Config) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, cg: cg}
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

	token, redirect, err := h.authUsecase.Login(req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.SetCookie("token", token, 3600*24, "/", "", false, true)

	response.Success(c, http.StatusOK, gin.H{
		"token":    token,
		"redirect": redirect,
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
