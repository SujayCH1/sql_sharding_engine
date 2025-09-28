package loader

import (
	"errors"
	"fmt"
	"sql_sharding_engine/internal/api"
	"sql_sharding_engine/internal/cache"
	"sql_sharding_engine/internal/repository/connections"
	"sql_sharding_engine/pkg/logger"

	"github.com/joho/godotenv"
)

// func to load all application services
func LoadServices() error {

	err := LoadEnv()
	if err != nil {
		return fmt.Errorf("failed to load application services: %w", err)
	}

	err = LoadAppAPIs()
	if err != nil {
		return fmt.Errorf("failed to load application services: %w", err)
	}

	err = LoadAppDB()
	if err != nil {
		return fmt.Errorf("failed to load application services: %w", err)
	}

	cache.CreateRedisClient()

	logger.Logger.Info("All application services loaded.")

	return nil
}

// func to laod environment variables
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return errors.New("failed to load env")
	}

	return nil
}

// funct expose all apis of application
func LoadAppAPIs() error {
	go func() {
		err := api.StartServer()
		if err != nil {
			logger.Logger.Error("failed to start server", "error", err)
		}
	}()

	return nil
}

// funct  to establish connection with application db
func LoadAppDB() error {
	err := connections.LoadMainDBConn()
	if err != nil {
		return fmt.Errorf("failed to load app DB: %w", err)
	}

	return nil
}
