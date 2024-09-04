package ucontext

import "github.com/gin-gonic/gin"

const (
	userIDKey   = "user_id_key_context"
	userNameKey = "user_name_key_context"
	userRoleKey = "user_role_key_context"
)

func GetUserID(c *gin.Context) string {
	val, ok := c.Get(userIDKey)
	if !ok {
		return ""
	}
	return val.(string)
}

func SetUserID(c *gin.Context, data interface{}) {
	c.Set(userIDKey, data)
}

func GetUserRole(c *gin.Context) string {
	val, ok := c.Get(userRoleKey)
	if !ok {
		return ""
	}
	return val.(string)
}

func SetUserRole(c *gin.Context, data interface{}) {
	c.Set(userRoleKey, data)
}

func GetUserName(c *gin.Context) string {
	val, ok := c.Get(userNameKey)
	if !ok {
		return ""
	}
	return val.(string)
}

func SetUserName(c *gin.Context, data interface{}) {
	c.Set(userRoleKey, data)
}
