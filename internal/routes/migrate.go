package routes

import (
	"goauth/internal/database"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func migrate(g *gin.RouterGroup) {
	g.POST("/migrate", func(c *gin.Context) {
		conn, ok := c.Get("DBConnection")
		if !ok {
			log.Fatal("Get db connection failed")
		}

		database.Migrate(conn.(*gorm.DB))
	})
}
