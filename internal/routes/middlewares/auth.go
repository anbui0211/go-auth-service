package routemiddlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	constants "goauth/internal/constant"
	"goauth/pkg/cache"
	"goauth/pkg/response"
	ujwt "goauth/utils/auth/jwt"
	ucontext "goauth/utils/context"

	"github.com/gin-gonic/gin"
)

// Authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authotization header is requested"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Verify token
		userId, userName, userRole, err := ujwt.VerifyTokenV2(tokenString, "access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		}

		ucontext.SetUserID(c, userId)
		ucontext.SetUserName(c, userName)
		ucontext.SetUserRole(c, userRole)
		c.Next()
	}
}

// Authorization
func RoleRequired(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := ucontext.GetUserRole(c)
		if userRole == "" {
			response.RMiddlewareError(c, http.StatusForbidden, "Access denied: User role is missing from the context.")
			return
		}

		// Check if the user no have a role
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.RMiddlewareError(c, http.StatusForbidden, "Insufficient permissions")
	}
}

// Check if user is banned
func CheckUserActive() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := ucontext.GetUserID(c)
		if userId == "" {
			response.RMiddlewareError(c, http.StatusForbidden, "Access denied: User ID is missing or invalid ")
			return
		}

		key := cache.GenKeyRedis("user_status", userId)
		userStatus, err := cache.GetRedis(context.Background(), key)
		if err != nil {
			response.RMiddlewareError(c, http.StatusInternalServerError, fmt.Sprintf("User status: %s", err.Error()))
			return
		}

		if userStatus == constants.StatusInActive {
			response.RMiddlewareError(c, http.StatusForbidden, "Access denied: User is banned")
		}

		c.Next()
	}
}
