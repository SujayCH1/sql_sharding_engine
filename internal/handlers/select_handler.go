package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sql_sharding_engine/internal/repository/connections"
	"sql_sharding_engine/internal/repository/database"
	"sql_sharding_engine/pkg/logger"
)

type DBHandler struct {
	Manager *database.CurrDBManager
}

// func to select db
func HandleSelectDB(m *database.CurrDBManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var db database.CurrDB
		if err := json.NewDecoder(r.Body).Decode(&db); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		m.SetCurrentDB(&db)

		if err := connections.ShardConnMgr.GetShardConnection(db.Name, db.ID); err != nil {
			logger.Logger.Error("Error while getting shard connections", "dbID", db.ID, "error", err)
			http.Error(w, "Shard Connection failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Current DB set to: %s", db.Name)
		logger.Logger.Info("Current DB updated", "db", db.Name)
	}
}
