package config

import (
	// "database/sql"
	sqlInst "database/sql"
	"time"

	"github.com/redis/go-redis/v9"
)

// Query struct for entire applications
type Query struct {
	QueryString string    `json:"query_string"`
	Timestamp   time.Time `json:"timestamp"`
	queryID     string    `json:"query_id"`
}

// Shard struct for entire application
type Shard struct {
	ShardName string `json:"shard_name"`
	ShardID   int    `json:"shard_id"`
	ShardHash uint32 `json:"shard_hash"`
	ShardHost string `json:"shard_host"`
	ShardPort int    `json:"shard_port"`
	ShardUser string `json:"shard_user"`
	ShardPass string `json:"shard_pass"`
}

// Database connection struct application database
type DBConnInfo struct {
	DBName string
	DBHost string
	DBPort int
	DBUser string
	DBPass string
}

// App DB connection credentials
var AppDBCommInfo *DBConnInfo

// App Db connection instance
var AppDBComm *sqlInst.DB

// redis client
var Redis *redis.Client

// temp pk of database
const KeyColumn string = "pk"
