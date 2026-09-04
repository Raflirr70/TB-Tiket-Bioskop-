package handler

import (
	"Project/internal/config"
	du "Project/internal/domain/usecase"
	"Project/pkg/response"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

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

func (h *FilmHandler) CreateFilm(c *gin.Context) {
	var req du.CreateFilmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name == "" {
		response.Error(c, http.StatusBadRequest, "nama film wajib diisi")
		return
	}

	film, err := h.filmUsecase.CreateFilm(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, film)
}

func (h *FilmHandler) UploadPoster(c *gin.Context) {
	file, err := c.FormFile("poster")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file poster tidak ditemukan")
		return
	}

	uploadDir := "./web/static/uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal membuat direktori upload")
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.Error(c, http.StatusInternalServerError, "gagal menyimpan file poster")
		return
	}

	fileURL := "/static/uploads/" + filename
	response.Success(c, http.StatusOK, gin.H{
		"url": fileURL,
	})
}
