package apis

import (
	"fmt"
	"net/http"
	"sql_sharding_engine/config"
	"sql_sharding_engine/services/database"
	"sql_sharding_engine/services/mapper"
	"sql_sharding_engine/services/parser"
	"time"
)

// func to expose backed on a local port
func StartServer() error {

	var CurrManager = &database.CurrDBManager{}

	mux := http.NewServeMux()

	mux.HandleFunc("/query", parser.HandleQuery)

	mux.HandleFunc("/shard", mapper.HandleShard)

	mux.HandleFunc("/database", database.HandleDatabase)

	mux.HandleFunc("/selectdb", CurrManager.HandleSelectDB)

	server := &http.Server{
		Addr:         ":8085",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	config.Logger.Info("Server listening at port 8085.")

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
