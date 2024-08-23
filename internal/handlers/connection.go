package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func connection(c *gin.Context) *gorm.DB {
	conn, ok := c.Get("DBConnection")
	if !ok {
		log.Fatal("Get db connection failed")
	}

	return conn.(*gorm.DB)
}

