package database

import (
	"fmt"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/internal/repository/database"
)

func FecthShards() ([]config.Shard, error) {
	currDB := database.CurrDBMgr.GetCurrentDB()

	query := fmt.Sprintf("SELECT * FROM %s", currDB.Name)

	shards, err := config.AppDBComm.Query(query)
	if err != nil {
		return []config.Shard{}, fmt.Errorf("Error: Failed to fetch Shards for view")
	}

	var Data []config.Shard

	for shards.Next() {
		var temp config.Shard

		if err := shards.Scan(&temp.ShardName, &temp.ShardID, &temp.ShardHash, &temp.ShardPort, &temp.ShardUser, &temp.ShardPass); err != nil {
			return []config.Shard{}, fmt.Errorf("failed to parrse rows for view : %s: %w", config.AppDBCommInfo.DBName, err)
		}

		Data = append(Data, temp)
	}

	return Data, nil

}
