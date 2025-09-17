package main

import (
	"log"

	"github.com/xuanviet96/seta-training/internal/cache"
	"github.com/xuanviet96/seta-training/internal/config"
	"github.com/xuanviet96/seta-training/internal/database"
	httpserver "github.com/xuanviet96/seta-training/internal/http"
	"github.com/xuanviet96/seta-training/internal/logger"
	"github.com/xuanviet96/seta-training/internal/search"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.New(cfg.AppEnv)

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL, logger)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Redis cache
	redis, err := cache.New(cfg, logger)
	if err != nil {
		log.Printf("Failed to initialize Redis: %v", err)
		// Continue without Redis if it fails
	}

	// Initialize Elasticsearch
	es, err := search.New(cfg, logger)
	if err != nil {
		log.Printf("Failed to initialize Elasticsearch: %v", err)
		// Continue without ES if it fails
	}

	// Ensure ES index exists
	if es != nil {
		if err := search.EnsureIndex(nil, es, cfg.ESIndex, logger); err != nil {
			log.Printf("Failed to ensure ES index: %v", err)
		}
	}

	// Initialize HTTP router
	router := httpserver.NewRouter(cfg, logger, db, redis, es)

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
