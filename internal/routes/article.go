package routes

import (
	"net/http"

	constants "goauth/internal/constant"
	routemiddlewares "goauth/internal/routes/middlewares"

	"github.com/gin-gonic/gin"
)

func article(r *gin.RouterGroup) {
	var g = r.Group("/articles")

	g.POST("/",
		routemiddlewares.AuthMiddleware(),
		routemiddlewares.CheckUserActive(),
		routemiddlewares.RoleRequired([]string{constants.RoleAdmin}),

		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "created successfully",
			})
		})
}
