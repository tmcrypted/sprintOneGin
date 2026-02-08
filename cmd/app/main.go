package main

import (
	"log"

	"sprin1/internal/config"
	"sprin1/internal/http"
	"sprin1/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	router := gin.Default()

	userService := service.NewUserService()
	http.RegisterRoutes(router, userService)

	addr := ":" + cfg.AppPort
	log.Printf("starting server on %s (env=%s)", addr, cfg.AppEnv)
	router.Run(addr)
}
