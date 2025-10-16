package database

import (
	"fmt"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/internal/repository/database"
)

func FetchShards() ([]config.Shard, error) {
	currDB := database.CurrDBMgr.GetCurrentDB()
	if currDB == nil {
		return nil, fmt.Errorf("current DB is not set")
	}

	if config.AppDBComm == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	query := fmt.Sprintf("SELECT * FROM %s", currDB.Name)

	rows, err := config.AppDBComm.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shards for view %s: %w", currDB.Name, err)
	}
	defer rows.Close()

	var data []config.Shard
	for rows.Next() {
		var temp config.Shard
		if err := rows.Scan(&temp.ShardName, &temp.ShardID, &temp.ShardHash, &temp.ShardHost, &temp.ShardPort, &temp.ShardUser, &temp.ShardPass); err != nil {
			return nil, fmt.Errorf("failed to parse rows for view %s: %w", currDB.Name, err)
		}
		data = append(data, temp)
	}

	return data, nil
}
