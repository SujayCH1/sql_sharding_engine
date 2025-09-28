package mapper

import (
	"fmt"
	"hash/crc32"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/pkg/logger"
	"strconv"
)

// func to handle shard add/ remove

// func used to add a new shard in mapping table
func AddShard(s config.Shard, tableName string) error {

	name := s.ShardName
	id := s.ShardID
	hash := CalcShardHash(strconv.Itoa(s.ShardID))
	host := s.ShardHost
	port := s.ShardPort
	user := s.ShardUser
	pass := s.ShardPass

	query := fmt.Sprintf(
		"INSERT INTO %s (database_name, shard_id, shard_hash, shard_host, shard_port, shard_user, shard_pass) VALUES (?, ?, ?, ?, ?, ?, ?)",
		tableName,
	)

	_, err := config.AppDBComm.Exec(query, name, id, hash, host, port, user, pass)
	if err != nil {
		return fmt.Errorf("failed to insert shard %s into %s: %w", s.ShardName, tableName, err)
	}

	if err != nil {
		return fmt.Errorf("failed to insert shard %s: %w", s.ShardName, err)
	}

	logger.Logger.Info("Added shard", name, "into mappings")

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

	logger.Logger.Info("Removed shard", s.ShardName, "from", tableName)
	return nil
}

// func to calc shard hash
func CalcShardHash(id string) uint32 {
	hash := crc32.ChecksumIEEE([]byte(id))

	return hash
}
