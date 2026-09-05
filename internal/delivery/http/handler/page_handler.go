package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageHandler struct{}

func NewPageHandler() *PageHandler {
	return &PageHandler{}
}

func (h *PageHandler) LandingPages(c *gin.Context) {
	email, _ := c.Get("email")
	firstname, _ := c.Get("firstname")
	c.HTML(http.StatusOK, "base", gin.H{
		"title":     "TB",
		"email":     email,
		"firstname": firstname,
		"page":      "landing",
		"nav":       false,
	})
}
func (h *PageHandler) HomePages(c *gin.Context) {
	email, _ := c.Get("email")
	firstname, _ := c.Get("firstname")
	c.HTML(http.StatusOK, "base", gin.H{
		"title":     "TB",
		"email":     email,
		"firstname": firstname,
		"page":      "home",
		"nav":       false,
	})
}

func (h *PageHandler) LoginPages(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "auth", gin.H{
		"title": "Login BooCinS",
		"mode":  "login",
		"nav":   false,
	})
}

func (h *PageHandler) RegisterPages(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "auth", gin.H{
		"title": "Register BooCinS",
		"mode":  "register",
		"nav":   false,
	})
}

func (h *PageHandler) DashboardPage(c *gin.Context) {
	email, _ := c.Get("email")
	firstname, _ := c.Get("firstname")
	c.HTML(http.StatusOK, "base", gin.H{
		"title":     "Dashboard BooCinS",
		"email":     email,
		"firstname": firstname,
		"page":      "dashboard",
		"nav":       true,
	})
}

func (h *PageHandler) DashboardSection(c *gin.Context, title string) {
	email, _ := c.Get("email")
	firstname, _ := c.Get("firstname")
	c.HTML(http.StatusOK, "base", gin.H{
		"title":     title + " - BooCinS",
		"email":     email,
		"firstname": firstname,
		"page":      "dashboard",
		"nav":       true,
		"section":   title,
	})
}

func (h *PageHandler) RoomsPage(c *gin.Context) {
	email, _ := c.Get("email")
	firstname, _ := c.Get("firstname")
	c.HTML(http.StatusOK, "base", gin.H{
		"title":     "Kelola Rooms & Seats - BooCinS",
		"email":     email,
		"firstname": firstname,
		"page":      "rooms",
		"nav":       true,
	})
}

func (h *PageHandler) ManageFilmsPage(c *gin.Context) {
	email, _ := c.Get("email")
	firstname, _ := c.Get("firstname")
	c.HTML(http.StatusOK, "base", gin.H{
		"title":     "Kelola Films & Schedules - BooCinS",
		"email":     email,
		"firstname": firstname,
		"page":      "manage-films",
		"nav":       true,
	})
}

func (h *PageHandler) CreateFilmPage(c *gin.Context) {
	email, _ := c.Get("email")
	firstname, _ := c.Get("firstname")
	c.HTML(http.StatusOK, "base", gin.H{
		"title":     "Tambah Film Baru - BooCinS",
		"email":     email,
		"firstname": firstname,
		"page":      "create-film",
		"nav":       true,
	})
}

func (h *PageHandler) FilmsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "base", gin.H{
		"title": "Daftar Movies",
		"page":  "films",
		"nav":   false,
	})
}
