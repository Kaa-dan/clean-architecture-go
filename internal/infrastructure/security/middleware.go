package security

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kaa-dan/clean-architecture-go/internal/domain/entities"
	"github.com/kaa-dan/clean-architecture-go/pkg/response"
)

type AuthMiddleware struct {
	jwtManager *JWTManager
}

func NewAuthMiddleware(jwtManager *JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := a.jwtManager.ValidateToken(tokenParts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func (a *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			response.Error(c, http.StatusUnauthorized, "User role not found")
			c.Abort()
			return
		}

		if role != string(entities.RoleAdmin) {
			response.Error(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
