package config

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// Query struct for entire applications
type Query struct {
	QueryString string    `json:"query_string"`
	Timestamp   time.Time `json:"timestamp"`
	queryID     string    `json:"query_id"`
}

// Shard struct for entire application
// Shard struct for entire application
type Shard struct {
	ShardName string `json:"shard_name"`
	ShardID   int    `json:"shard_id"`
	ShardHash uint32 `json:"shard_hash"`
	ShardHost string `json:"shard_host"`
	ShardPort int    `json:"shard_port"`
}

// Database connection struct for entire application
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
var AppDBComm *sql.DB

// services logger
var Logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// redis client
var Redis *redis.Client

// temp pk of database
const KeyColumn string = "pk"
