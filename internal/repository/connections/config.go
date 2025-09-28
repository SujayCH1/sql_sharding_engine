package connections

import (
	"database/sql"
	"sync"
)

// to handle connection of all sahrds for current selected database
type ActiveDBShardConnectionManager struct {
	mu         sync.Mutex
	ActiveDB   int
	ShardConns map[int]map[int]*sql.DB
}

// Current establish shard connections
var ShardConnMgr = &ActiveDBShardConnectionManager{}

type CurrDBInfo struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type CurrDBGetter interface {
	GetID() int
	GetName() string
}
