package database

import (
	"database/sql"
	"sync"
)

type Database struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type DBReq struct {
	ReqType string   `json:"type"`
	DBInfo  Database `json:"database"`
}

type CurrDB struct {
	Name string
	ID   int
}

type CurrDBManager struct {
	mu         sync.Mutex
	curr       *CurrDB
	shardConns map[string]*sql.DB
}

// Current selected databse instance
var CurrDBMgr = &CurrDBManager{}

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
