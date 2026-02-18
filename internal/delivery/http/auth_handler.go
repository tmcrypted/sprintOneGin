package http

import (
	"net/http"
	"time"

	"sprin1/internal/delivery/http/dto"

	"github.com/gin-gonic/gin"
)

const refreshCookieName = "refresh_token"

// Cookie TTL для refresh-токена должен быть не меньше, чем refreshTokenTTL в сервисе.
const refreshCookieTTL = 30 * 24 * time.Hour

func (s *Server) RegisterAuthRoutes() {
	g := s.router.Group("/auth")
	g.POST("/login", s.login())
	g.POST("/register", s.register())
	g.POST("/refresh", s.refresh())
}

func (s *Server) login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.LoginRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := s.authService.Login(c.Request.Context(), body)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		setRefreshCookie(c, resp.RefreshToken)
		c.JSON(http.StatusOK, resp)
	}
}

func (s *Server) register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.CreateUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := s.authService.Register(c.Request.Context(), body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		setRefreshCookie(c, resp.RefreshToken)
		c.JSON(http.StatusCreated, resp)
	}
}

func (s *Server) refresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(refreshCookieName)
		if err != nil || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
			return
		}

		resp, err := s.authService.Refresh(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		setRefreshCookie(c, resp.RefreshToken)
		c.JSON(http.StatusOK, resp)
	}
}

func setRefreshCookie(c *gin.Context, token string) {
	// HttpOnly, Secure=false (для локальной разработки), Path="/"
	maxAge := int(refreshCookieTTL.Seconds())
	c.SetCookie(refreshCookieName, token, maxAge, "/", "", false, true)
}

