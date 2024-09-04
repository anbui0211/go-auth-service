package routes

import (
	"goauth/internal/dao"
	"goauth/internal/handlers"
	"goauth/internal/services"

	"github.com/gin-gonic/gin"
)

func user(r *gin.RouterGroup) {
	var (
		g = r.Group("/users")
		d = dao.NewUserDao()
		s = services.NewuserService(d)
		h = handlers.NewUserHandler(s)
	)

	g.PATCH("/:id/status", h.ChangeStatus)
}