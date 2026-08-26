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
	c.HTML(http.StatusOK, "base", gin.H{
		"title": "TB",
		"email": email,
		"page":  "landing",
	})
}
