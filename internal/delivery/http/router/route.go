package router

import (
	"Project/internal/config"
	"Project/internal/delivery/http/handler"
	"Project/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	filmHandler *handler.FilmHandler,
	roomHandler *handler.RoomHandler,
) *gin.Engine {
	r := gin.Default()

	r.POST("/api/v1/auth/login", authHandler.Login)
	r.POST("/api/v1/auth/logout", authHandler.Logout)
	r.POST("/api/v1/auth/register", authHandler.Register)

	r.GET("/api/v1/films", filmHandler.GetAllFilm)
	r.POST("/api/v1/films", middleware.RequireAdmin(cfg.JWT), filmHandler.CreateFilm)
	r.POST("/api/v1/films/upload", middleware.RequireAdmin(cfg.JWT), filmHandler.UploadPoster)

	// Rooms & Seats
	r.GET("/api/v1/rooms", middleware.RequireAdmin(cfg.JWT), roomHandler.GetAllRooms)
	r.POST("/api/v1/rooms", middleware.RequireAdmin(cfg.JWT), roomHandler.CreateRoom)
	r.DELETE("/api/v1/rooms/:id", middleware.RequireAdmin(cfg.JWT), roomHandler.DeleteRoom)
	r.PUT("/api/v1/rooms/:id", middleware.RequireAdmin(cfg.JWT), roomHandler.UpdateRoom)

	return r
}
