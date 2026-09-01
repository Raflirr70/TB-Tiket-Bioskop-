package handler

import (
	"Project/internal/config"
	du "Project/internal/domain/usecase"
	"Project/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FilmHandler struct {
	filmUsecase du.FilmUseCase
	cg          *config.Config
}

func NewFilmHandler(filmUsecase du.FilmUseCase, cg *config.Config) *FilmHandler {
	return &FilmHandler{filmUsecase: filmUsecase, cg: cg}
}

func (h *FilmHandler) GetAllFilm(c *gin.Context) {
	films, err := h.filmUsecase.GetAllFilm()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, films)
}
