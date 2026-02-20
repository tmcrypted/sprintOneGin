package http

import (
	"net/http"
	"strings"

	"sprin1/internal/model"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware проверяет JWT access-токен из заголовка Authorization: Bearer <token>
// и кладёт пользователя в контекст Gin под ключом "user".
func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		tokenStr := parts[1]

		user, err := s.authService.ParseToken(c.Request.Context(), tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// ModeratorMiddleware проверяет, что аутентифицированный пользователь имеет роль модератора.
// Должен использоваться после AuthMiddleware.
func (s *Server) ModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
			return
		}

		user, ok := userAny.(*model.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user type in context"})
			return
		}

		if user.Role != model.RoleModerator {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "moderator access required"})
			return
		}

		c.Next()
	}
}

