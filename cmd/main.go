package main

import (
	"fmt"
	"log"

	"goauth/internal/server"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("internal/templates/*")

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("error loading environment")
	}
	fmt.Println("Load env success ...")

	server.Init(r)

	r.Run(":8000")
}
