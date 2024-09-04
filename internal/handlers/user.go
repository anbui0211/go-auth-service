package handlers

import (
	"goauth/internal/services"

	"github.com/gin-gonic/gin"
)

type IUserhandler interface {
	ChangeStatus(c *gin.Context)
}

type userHandler struct {
	userSvc services.IUserService
}

func NewUserHandler(userSvc services.IUserService) IUserhandler {
	return userHandler{
		userSvc: userSvc,
	}
}
func (uh userHandler) ChangeStatus(c *gin.Context) {
	var (
		userId = c.Param("id")
		conn   = connection(c)
	)

	if userId == "" {
		c.JSON(400, gin.H{"error": "user id invalid"})
		return
	}

	if err := uh.userSvc.ChangeStatus(conn, userId); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "change user status success"})
}
