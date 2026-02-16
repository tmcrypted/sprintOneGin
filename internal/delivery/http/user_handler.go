package http

import (
	"net/http"
	"sprin1/internal/delivery/http/dto"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterUserRoutes() {
	g := s.router.Group("/users")
	g.POST("/create", s.createUser())
	g.GET("/:id", s.getUser())
	g.GET("/", s.getAllUsers())
}

func (s *Server) createUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.CreateUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := s.userService.CreateUser(c.Request.Context(), body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func (s *Server) getUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		user, err := s.userService.GetUser(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func (s *Server) getAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.userService.GetAllUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}
