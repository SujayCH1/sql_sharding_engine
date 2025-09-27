package mapper

import "sql_sharding_engine/config"

type shardReq struct {
	ReqType string          `json:"type"`
	Shard   config.Shard    `json:"shard"`
	DB      config.Database `json:"database"`
}
