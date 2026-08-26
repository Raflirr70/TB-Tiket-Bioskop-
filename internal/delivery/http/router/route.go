package router

import (
	"Project/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func UserRoute(route *gin.Engine, userHandler *handler.UserHandler) {
	route.LoadHTMLGlob("web/templates/**/*")

}
