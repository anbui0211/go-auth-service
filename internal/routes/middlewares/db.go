package routemiddlewares

import (
	"goauth/internal/database"

	"github.com/gin-gonic/gin"
)

func OpenConnection() gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := database.Connect()
		c.Set("DBConnection", conn)
		c.Next()
	}
}
