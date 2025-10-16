package parser

import (
	"fmt"
	"sql_sharding_engine/internal/repository/connections"
	database "sql_sharding_engine/internal/repository/database/controller"
	"sql_sharding_engine/pkg/logger"

	"github.com/xwb1989/sqlparser"
)

func ExecuteDDL(stmt sqlparser.Statement, queryString string) error {
	logger.Logger.Info("DDL query received", "query", queryString)

	// Fetch all shards for the currently selected DB
	shards, err := database.FetchShards()
	if err != nil {
		return fmt.Errorf("failed to fetch shards: %w", err)
	}

	if len(shards) == 0 {
		return fmt.Errorf("no shards found for current database")
	}

	mgr := connections.ShardConnMgr
	dbID := mgr.ActiveDB

	if dbID == 0 {
		return fmt.Errorf("no active database selected")
	}

	mgr.Mu.Lock()
	defer mgr.Mu.Unlock()

	// Iterate through all shard connections for the current DB
	for _, shard := range shards {
		shardID := shard.ShardID
		conn, ok := mgr.ShardConns[dbID][shardID]

		if !ok || conn == nil {
			logger.Logger.Warn("Skipping nil or missing connection for shard", "shardID", shardID)
			continue
		}

		dbName := shard.ShardName

		// Switch to the shard's database
		useQuery := fmt.Sprintf("USE %s;", dbName)
		if _, err := conn.Exec(useQuery); err != nil {
			logger.Logger.Error("Failed to select database on shard", "shardID", shardID, "error", err)
			return fmt.Errorf("failed to select database on shard %d: %w", shardID, err)
		}

		// Execute the DDL query
		if _, err := conn.Exec(queryString); err != nil {
			logger.Logger.Error("Failed to execute DDL on shard", "shardID", shardID, "error", err)
			return fmt.Errorf("failed on shard %d: %w", shardID, err)
		}

		logger.Logger.Info("✅ DDL executed successfully on shard", "shardID", shardID, "dbName", dbName)
	}

	return nil
}

// executeDML routes DML query to the correct shard based on primary key hash
func executeDML(stmt sqlparser.Statement, keys []string, queryString string) error {
	if len(keys) == 0 {
		return fmt.Errorf("no primary keys provided to route DML")
	}

	mgr := connections.ShardConnMgr
	dbID := mgr.ActiveDB
	if dbID == 0 {
		return fmt.Errorf("no active database selected")
	}

	mgr.Mu.Lock()
	shards, ok := mgr.ShardConns[dbID]
	mgr.Mu.Unlock()

	if !ok || len(shards) == 0 {
		return fmt.Errorf("no shard connections found for DB %d", dbID)
	}

	for _, key := range keys {
		shardID, err := getShardIDForKey(key, shards)
		if err != nil {
			return fmt.Errorf("failed to get shard for key %s: %w", key, err)
		}

		conn, ok := shards[shardID]
		if !ok || conn == nil {
			return fmt.Errorf("no connection found for shard %d", shardID)
		}

		_, err = conn.Exec(queryString)
		if err != nil {
			logger.Logger.Error("Failed to execute DML on shard", "shardID", shardID, "error", err)
			return fmt.Errorf("failed to execute on shard %d: %w", shardID, err)
		}

		logger.Logger.Info("DML executed successfully on shard", "shardID", shardID, "key", key)
	}

	return nil
}
