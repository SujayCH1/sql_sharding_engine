package database

import (
	"database/sql"
	"sql_sharding_engine/config"
	"sync"
)

type dbReq struct {
	ReqType string          `json:"type"`
	DBInfo  config.Database `json:"database"`
}

type CurrDB struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type CurrDBManager struct {
	mu         sync.Mutex
	curr       *CurrDB
	shardConns map[string]*sql.DB
}

// Thread-safe setter
func (m *CurrDBManager) SetCurrentDB(db *CurrDB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.curr = db
}

// Thread-safe getter
func (m *CurrDBManager) GetCurrentDB() *CurrDB {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.curr
}
