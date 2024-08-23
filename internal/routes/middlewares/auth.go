package routemiddlewares

import (
	"net/http"
	"strings"

	"goauth/pkg/response"
	ujwt "goauth/utils/auth/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenKeyContext = "TOKEN_KEY_CONTEXT"
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
		token, err := ujwt.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
		}

		c.Set(tokenKeyContext, token)
		c.Next()
	}
}

// Authorization
func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := c.Get(tokenKeyContext)
		if !ok {
			response.RMiddlewareError(c, http.StatusUnauthorized, "Token does not exist in context")
			return
		}

		claims, ok := token.(*jwt.Token).Claims.(jwt.MapClaims)
		if !ok {
			response.RMiddlewareError(c, http.StatusUnauthorized, "invalid token claims")
			return
		}

		userRole, ok := claims["role"].(string)
		if !ok {
			response.RMiddlewareError(c, http.StatusUnauthorized, "role not found in token")
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
