package handler

import (
	"Project/internal/config"
	du "Project/internal/domain/usecase"
	"Project/pkg/response"
	"net/http"
	"strconv"

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
	limit := 0

	if limitParam := c.Query("limit"); limitParam != "" {
		var err error

		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit < 0 {
			response.Error(c, http.StatusBadRequest, "limit harus berupa angka")
			return
		}
	}

	sort := c.Query("sort")

	films, err := h.filmUsecase.GetAllFilm(limit, sort)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, films)
}
