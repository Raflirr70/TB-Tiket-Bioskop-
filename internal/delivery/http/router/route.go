package router

import (
	"Project/internal/config"
	"Project/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(
	cfg *config.Config,
	pageHandler *handler.PageHandler,
) *gin.Engine {
	r := gin.Default()

	// Load HTML templates and serve static files
	r.LoadHTMLGlob("web/templates/**/*")
	r.Static("/static", "./web/static")

	// Public Web Routes
	// r.GET("/", middleware.Auth(cfg), pageHandler.GetLandingPage) // we use Auth middleware but it won't force-abort page request if cookie isn't present unless we enforce it in Handler. Actually landing page is public, but let's allow it to read auth info
	r.GET("/", pageHandler.LandingPages)
	return r
}
