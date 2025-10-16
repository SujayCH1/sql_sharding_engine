package handlers

import (
	"encoding/json"
	"net/http"
	"sql_sharding_engine/internal/config"
	database "sql_sharding_engine/internal/repository/database/controller"
	"sql_sharding_engine/pkg/logger"
)

func HandleDBMappingView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	datbases, err := database.FetchDBMappings()
	if err != nil {
		http.Error(w, "Error: Internal Server Error", http.StatusInternalServerError)
		logger.Logger.Error("Error: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(datbases)
}

func HandleShardView(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}

	shards, err := database.FecthShards()
	if err != nil {
		http.Error(w, "Error: Internal Server Error", http.StatusInternalServerError)
		logger.Logger.Error("Error: %s", err)
		return
	}

	if len(shards) == 0 {
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode([]config.Shard{})
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(shards)

}
