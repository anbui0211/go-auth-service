package routes

import (
	"goauth/internal/dao"
	"goauth/internal/handlers"
	"goauth/internal/services"

	"github.com/gin-gonic/gin"
)

func auth(g *gin.RouterGroup) {
	var (
		r = g.Group("/auth")

		d = dao.NewUserDao()
		s = services.NewAuthService(d)
		h = handlers.NewAuthHandler(s)
	)

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.POST("/refresh", h.Refresh)

	// OAuth2 authentication
	r.GET("/", h.MainAuthPage)
	r.GET("/oauth-login", h.StartGoogleOAuth)
	r.GET("/callback", h.OAuthGoogleCallback)
}
