package mapper

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"net/http"
	"sql_sharding_engine/config"
	"strconv"
)

type shardReq struct {
	ReqType string          `json:"type"`
	Shard   config.Shard    `json:"shard"`
	DB      config.Database `json:"database"`
}

// func to handle shard add/ remove
func HandleShard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req shardReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	switch req.ReqType {
	case "add":
		err := AddShard(req.Shard, req.DB.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "remove":
		err := RemoveShard(req.Shard, req.DB.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "invalid type", http.StatusBadRequest)
	}

}

// func used to add a new shard in mapping table
func AddShard(s config.Shard, tableName string) error {

	name := s.ShardName
	id := s.ShardID
	hash := CalcShardHash(strconv.Itoa(s.ShardID))
	host := s.ShardHost
	port := s.ShardPort

	query := fmt.Sprintf(
		"INSERT INTO %s (database_name, shard_id, shard_hash, shard_host, shard_port) VALUES (?, ?, ?, ?, ?)",
		tableName,
	)

	_, err := config.AppDBComm.Exec(query, name, id, hash, host, port)
	if err != nil {
		return fmt.Errorf("failed to insert shard %s into %s: %w", s.ShardName, tableName, err)
	}

	if err != nil {
		return fmt.Errorf("failed to insert shard %s: %w", s.ShardName, err)
	}

	config.Logger.Info("Added shard", name, "into mappings")

	return nil
}

// func to remove a existing shard from mapping table
func RemoveShard(s config.Shard, tableName string) error {
	id := s.ShardID

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE shard_id = ?",
		tableName,
	)

	_, err := config.AppDBComm.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to remove shard %s from %s: %w", s.ShardName, tableName, err)
	}

	config.Logger.Info("Removed shard", s.ShardName, "from", tableName)
	return nil
}

// func to calc shard hash
func CalcShardHash(id string) uint32 {
	hash := crc32.ChecksumIEEE([]byte(id))

	return hash
}
