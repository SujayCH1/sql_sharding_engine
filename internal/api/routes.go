package api

import (
	"fmt"
	"net/http"
	"sql_sharding_engine/internal/handlers"
	"sql_sharding_engine/internal/repository/database"
	"sql_sharding_engine/pkg/logger"
	"time"
)

// func to expose backed on a local port
func StartServer() error {

	var CurrManager = &database.CurrDBManager{}

	mux := http.NewServeMux()

	mux.HandleFunc("/query", handlers.HandleQuery)

	mux.HandleFunc("/shard", handlers.HandleShard)

	mux.HandleFunc("/database", handlers.HandleDatabase)

	mux.HandleFunc("/selectdb", handlers.HandleSelectDB(CurrManager))

	server := &http.Server{
		Addr:         ":8085",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Logger.Info("Server listening at port 8085.")

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
