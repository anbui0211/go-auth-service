package main

import (
	"goauth/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading environment")
	}

	r.LoadHTMLGlob("internal/templates/*")
	routes.Init(r)

	r.Run(":8000")
}
