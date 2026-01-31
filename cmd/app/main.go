package main

import (
	"sprin1/internal/http"
	"sprin1/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	userService := service.NewUserService()
	http.RegisterRoutes(router, userService)

	router.Run(":8080")
}
