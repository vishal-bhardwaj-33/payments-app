package database

import (
	"fmt"
	"log"

	"payments-app/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection
func InitDB(cfg config.Config) error {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host, cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.Port, cfg.Database.SslMode,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}

	log.Println("Database initialized successfully.")
	return nil
}
