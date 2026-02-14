package http

import (
	"net/http"
	"strconv"

	"sprin1/internal/model"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, userService UserService) {
	g := r.Group("/users")
	g.POST("/create", createUser(userService))
	g.GET("/:id", getUser(userService))
	g.GET("/", getAllUsers(userService))
}

func createUser(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email    string         `json:"email" binding:"required"`
			Password string         `json:"password" binding:"required"`
			FIO      string         `json:"fio" binding:"required"`
			Role     model.UserRole `json:"role"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := userService.CreateUser(c.Request.Context(), body.Email, body.Password, body.FIO, body.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func getUser(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		user, err := userService.GetUser(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func getAllUsers(userService UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userService.GetAllUsers(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}
