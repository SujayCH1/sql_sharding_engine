package connections

import (
	"database/sql"
	"sync"
)

// to handle connection of all sahrds for current selected database
type ActiveDBShardConnectionManager struct {
	mu         sync.Mutex
	ActiveDB   int
	ShardConns map[string]*sql.DB
}

func (m *ActiveDBShardConnectionManager) SetShardConn(name string, conn *sql.DB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.ShardConns == nil {
		m.ShardConns = make(map[string]*sql.DB)
	}
	m.ShardConns[name] = conn
}

func (m *ActiveDBShardConnectionManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, conn := range m.ShardConns {
		conn.Close()
	}
	m.ShardConns = make(map[string]*sql.DB)
}
