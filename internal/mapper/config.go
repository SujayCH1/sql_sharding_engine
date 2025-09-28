package mapper

import (
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/internal/repository/database"
)

type ShardReq struct {
	ReqType string            `json:"type"`
	Shard   config.Shard      `json:"shard"`
	DB      database.Database `json:"database"`
}
