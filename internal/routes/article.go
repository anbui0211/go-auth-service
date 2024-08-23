package routes

import (
	"net/http"

	routemiddlewares "goauth/internal/routes/middlewares"

	"github.com/gin-gonic/gin"
)

func article(r *gin.RouterGroup) {
	var g = r.Group("/articles")

	g.POST("/",
		routemiddlewares.AuthMiddleware(),
		routemiddlewares.RoleRequired("ADMIN"),

		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "created successfully",
			})
		})
}
