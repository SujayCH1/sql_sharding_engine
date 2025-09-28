package mapper

import (
	"sql_sharding_engine/config"
	"sql_sharding_engine/services/database"
)

type shardReq struct {
	ReqType string            `json:"type"`
	Shard   config.Shard      `json:"shard"`
	DB      database.Database `json:"database"`
}
