package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string `env:"APP_PORT"`
	AppEnv  string `env:"APP_ENV"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASSWORD"`
	DBName     string `env:"DB_NAME"`
	DBSSLMode  string `env:"DB_SSLMODE"`

	JWTSecret         string `env:"JWT_SECRET"`
	JWTTTLMinutes     int    `env:"JWT_TTL_MINUTES" envDefault:"60"`
	JWTREFRESHTTLDays int    `env:"JWT_REFRESH_TTL_DAYS" envDefault:"60"`

	// KafkaBrokers — список адресов брокеров Kafka через запятую, например: "kafka:9092".
	KafkaBrokers string `env:"KAFKA_BROKERS"`
	// KafkaRatingTopic — имя топика для событий пересчёта рейтинга пользователя.
	KafkaRatingTopic string `env:"KAFKA_RATING_TOPIC" envDefault:"user-rating-recalc"`
	// KafkaRatingGroupID — consumer group id для пересчёта рейтинга.
	KafkaRatingGroupID string `env:"KAFKA_RATING_GROUP_ID" envDefault:"sprin1-rating-consumer"`
}

// Load загружает .env (если есть) и заполняет Config из переменных окружения по тегам. Без дефолтов.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *Config) DSN() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode)
	fmt.Println(dsn)
	return dsn
}

// AccessTTL возвращает время жизни access токена.
func (c *Config) AccessTTL() time.Duration {
	return time.Duration(c.JWTTTLMinutes) * time.Minute
}

// RefreshTTL возвращает время жизни refresh токена.
func (c *Config) RefreshTTL() time.Duration {
	return time.Duration(c.JWTREFRESHTTLDays) * 24 * time.Hour
}
