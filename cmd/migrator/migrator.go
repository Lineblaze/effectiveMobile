package migrator

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

var dbConnString string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbConnString = os.Getenv("POSTGRES_CONN")
	if dbConnString == "" {
		log.Fatal("POSTGRES_CONN is not set in the .env file")
	}
}

func Migrate() {
	m, err := migrate.New("file://migration", dbConnString)
	if err != nil {
		log.Fatalf("Error initializing migration: %v", err)
	}

	migrationErr := m.Up()
	if migrationErr != nil {
		if errors.Is(migrationErr, migrate.ErrNoChange) {
			fmt.Println("Migration: No changes")
		} else {
			log.Fatalf("Migration error: %v", migrationErr)
		}
	} else {
		fmt.Println("Database migration completed successfully.")
	}
}
