package http

import (
	"net/http"
	"strconv"

	"sprin1/internal/service"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userService service.UserService) {
	r.POST("/createUser", createUserHandler(userService))
	r.GET("/users/:id", getUserHandler(userService))
	r.GET("/users", getAllUsersHandler(userService))
}

func createUserHandler(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Name string `json:"name" binding:"required"`
			Age  int    `json:"age"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := userService.CreateUser(body.Name, body.Age)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func getUserHandler(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}
		user, err := userService.GetUser(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func getAllUsersHandler(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users := userService.GetAllUsers()
		c.JSON(http.StatusOK, users)
	}
}
