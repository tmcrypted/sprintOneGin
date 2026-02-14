package main

import (
	"context"
	"log"

	"sprin1/internal/config"
	"sprin1/internal/delivery/http"
	"sprin1/internal/repository/postgres"
	"sprin1/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	userRepo := postgres.NewUserRepository(pool)
	userService := service.NewUserService(userRepo)

	reviewRepo := postgres.NewReviewRepository(pool)
	reviewService := service.NewReviewService(reviewRepo, userRepo)

	router := gin.Default()
	http.RegisterRoutes(router, userService, reviewService)

	addr := ":" + cfg.AppPort
	log.Printf("starting server on %s (env=%s)", addr, cfg.AppEnv)
	router.Run(addr)
}
