package main

import (
	"effectiveMobile/config"
	"effectiveMobile/internal/httpServer"
	"effectiveMobile/pkg/logger"
	"log"
)

func main() {
	log.Println("Starting server")

	cfg := config.LoadConfig()
	log.Println("Config loaded")

	appLogger := logger.NewApiLogger(cfg)
	err := appLogger.InitLogger()
	if err != nil {
		log.Fatalf("Cannot init logger: %v", err.Error())
	}
	log.Println("Logger loaded")

	s := httpServer.NewServer(cfg, appLogger)
	if err = s.Run(); err != nil {
		appLogger.Errorf("Server run error: %v", err)
	}
}
