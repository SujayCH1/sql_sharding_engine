package handlers

import (
	"encoding/json"
	"net/http"
	"sql_sharding_engine/internal/repository/database"
	"sql_sharding_engine/pkg/logger"
)

// func to add/remove database
func HandleDatabase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newReq database.DBReq

	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	switch newReq.ReqType {
	case "add":
		err := database.AddDatabase(newReq.DBInfo)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			logger.Logger.Error("failed to add a new database", err)
		}

	case "delete":
		err := database.DeleteDatabase(newReq.DBInfo)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			logger.Logger.Error("failed to delete exsisting database:", err)
		}

	default:
		http.Error(w, "invalid type", http.StatusBadRequest)
	}

}
