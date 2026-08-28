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

func (h *PageHandler) LoginPages(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "Login BooCinS",
		"nav":   false,
	})
}

func (h *PageHandler) RegisterPages(c *gin.Context) {
	if _, err := c.Cookie("token"); err == nil {
		c.Redirect(http.StatusSeeOther, "/dashboard")
		return
	}
	c.HTML(http.StatusOK, "register", gin.H{
		"title": "Register BooCinS",
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
