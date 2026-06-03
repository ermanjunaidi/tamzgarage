package config

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	AllowedOrigins string
}

func Load() *Config {
	godotenv.Load("../.env", ".env")

	return &Config{
		Port:           getEnv("PORT", "8082"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://bengkelpro:bengkelpro123@localhost:5432/bengkelpro?sslmode=disable"),
		JWTSecret:      getEnv("JWT_SECRET", "bengkelpro-secret-key-2024"),
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
