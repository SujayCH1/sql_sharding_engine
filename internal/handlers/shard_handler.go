package handlers

import (
	"encoding/json"
	"net/http"
	"sql_sharding_engine/internal/mapper"
)

func HandleShard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req mapper.ShardReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	switch req.ReqType {
	case "add":
		err := mapper.AddShard(req.Shard, req.DB.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "delete":
		err := mapper.RemoveShard(req.Shard, req.DB.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "invalid type", http.StatusBadRequest)
	}

}
