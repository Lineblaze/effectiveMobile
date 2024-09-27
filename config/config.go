package config

import (
	"log"
	"os"
)

type Config struct {
	ServiceName string

	Postgres struct {
		ConnURL  string
		JDBCURL  string
		User     string
		Password string
		Host     string
		Port     string
		DBName   string
		SSLMode  string
		PgDriver string
	}

	Server struct {
		Address                     string
		ShowUnknownErrorsInResponse bool
	}
}

func LoadConfig() *Config {
	c := &Config{
		ServiceName: "Tenders",
		Postgres: struct {
			ConnURL  string
			JDBCURL  string
			User     string
			Password string
			Host     string
			Port     string
			DBName   string
			SSLMode  string
			PgDriver string
		}{
			ConnURL:  os.Getenv("POSTGRES_CONN"),
			JDBCURL:  os.Getenv("POSTGRES_JDBC_URL"),
			User:     os.Getenv("POSTGRES_USERNAME"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			DBName:   os.Getenv("POSTGRES_DATABASE"),
			SSLMode:  "disable",
			PgDriver: "pgx",
		},
		Server: struct {
			Address                     string
			ShowUnknownErrorsInResponse bool
		}{
			Address:                     os.Getenv("SERVER_ADDRESS"),
			ShowUnknownErrorsInResponse: false,
		},
	}

	if c.Postgres.ConnURL == "" || c.Server.Address == "" {
		log.Fatalf("Required environment variables not set")
	}

	return c
}
