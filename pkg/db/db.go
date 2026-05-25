package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		envOrDefault("DB_HOST", "localhost"),
		envOrDefault("DB_USER", "postgres"),
		os.Getenv("DB_PASSWORD"),
		envOrDefault("DB_NAME", "postgres"),
		envOrDefault("DB_PORT", "5432"),
		envOrDefault("DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db
}

func envOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
