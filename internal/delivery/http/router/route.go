package router

import (
	"Project/internal/config"
	"Project/internal/delivery/http/handler"
	"Project/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	pageHandler *handler.PageHandler,
	authHandler *handler.AuthHandler,
	filmHandler *handler.FilmHandler,
) *gin.Engine {
	r := gin.Default()

	// Load HTML templates and serve static files
	r.LoadHTMLGlob("web/templates/**/*")
	r.Static("/static", "./web/static")

	// Public Web Routes
	// r.GET("/", middleware.Auth(cfg), pageHandler.GetLandingPage) // we use Auth middleware but it won't force-abort page request if cookie isn't present unless we enforce it in Handler. Actually landing page is public, but let's allow it to read auth info
	r.GET("/", middleware.OptionalAuth(cfg.JWT), pageHandler.LandingPages)
	r.GET("/home", middleware.OptionalAuth(cfg.JWT), pageHandler.HomePages)
	r.GET("/films", middleware.OptionalAuth(cfg.JWT), pageHandler.FilmsPage)
	r.GET("/login", middleware.OptionalAuth(cfg.JWT), pageHandler.LoginPages)
	r.GET("/register", middleware.OptionalAuth(cfg.JWT), pageHandler.RegisterPages)

	r.POST("/api/v1/auth/login", authHandler.Login)
	r.POST("/api/v1/auth/logout", authHandler.Logout)
	r.POST("/api/v1/auth/register", authHandler.Register)

	r.GET("/api/v1/films", filmHandler.GetAllFilm)

	r.GET("/dashboard", middleware.RequireAdmin(cfg.JWT), pageHandler.DashboardPage)

	// Dashboard sections (placeholder, judul saja)
	r.GET("/transactions", middleware.RequireAdmin(cfg.JWT), func(c *gin.Context) { pageHandler.DashboardSection(c, "Transaction Ticket Offline") })
	r.GET("/schedules", middleware.RequireAdmin(cfg.JWT), func(c *gin.Context) { pageHandler.DashboardSection(c, "Daftar Schedule Film") })
	r.GET("/manage-films", middleware.RequireAdmin(cfg.JWT), func(c *gin.Context) { pageHandler.DashboardSection(c, "Kelola Films & Schedules") })
	r.GET("/manage-rooms", middleware.RequireAdmin(cfg.JWT), func(c *gin.Context) { pageHandler.DashboardSection(c, "Kelola Rooms & Seats") })
	r.GET("/food", middleware.RequireAdmin(cfg.JWT), func(c *gin.Context) { pageHandler.DashboardSection(c, "Kelola Makanan") })
	r.GET("/membership", middleware.RequireAdmin(cfg.JWT), func(c *gin.Context) { pageHandler.DashboardSection(c, "Kelola Membership") })

	return r
}
