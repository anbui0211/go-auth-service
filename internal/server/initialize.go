package server

import (
	"goauth/internal/routes"
	"goauth/pkg/cache"
	zlog "goauth/pkg/log"

	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	// Package
	zlog.Init("auth-service")
	cache.ConnectRedis()

	routes.Init(e)
}
