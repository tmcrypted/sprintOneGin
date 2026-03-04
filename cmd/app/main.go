package main

import (
	"context"
	"log"
	"time"

	"sprin1/internal/config"
	"sprin1/internal/delivery/http"
	"sprin1/internal/repository/postgres"
	"sprin1/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

	pvzRepo := postgres.NewPVZRepository(pool)
	pvzService := service.NewPVZService(pvzRepo)

	refreshRepo := postgres.NewRefreshSessionRepository(pool)
	authService := service.NewAuthService(userRepo, refreshRepo, cfg.JWTSecret, cfg.AccessTTL(), cfg.RefreshTTL())

	srv := http.NewServer(userService, reviewService, authService, pvzService, cfg.JWTSecret)

	addr := ":" + cfg.AppPort
	log.Printf("starting server on %s (env=%s)", addr, cfg.AppEnv)
	if err := srv.Run(addr); err != nil {
		log.Fatalln("Не удалось запустить сервер")
	}
}
