package routes

import (
	routemiddlewares "goauth/internal/routes/middlewares"

	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	var g = e.Group("")

	// Set connection
	g.Use(routemiddlewares.OpenConnection())

	migrate(g)
	auth(g)
	article(g)
}
