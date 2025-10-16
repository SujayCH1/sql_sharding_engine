package connections

import (
	"database/sql"
	"fmt"
	"sql_sharding_engine/internal/config"
	// "sql_sharding_engine/services/database"
)

func (m *ActiveDBShardConnectionManager) SetShardConn(dbID, shardID int, conn *sql.DB) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if m.ShardConns == nil {
		m.ShardConns = make(map[int]map[int]*sql.DB)
	}
	if m.ShardConns[dbID] == nil {
		m.ShardConns[dbID] = make(map[int]*sql.DB)
	}
	m.ShardConns[dbID][shardID] = conn
}

func (m *ActiveDBShardConnectionManager) GetShardConn(dbID, shardID int) (*sql.DB, bool) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	shardMap, ok := m.ShardConns[dbID]
	if !ok {
		return nil, false
	}
	conn, ok := shardMap[shardID]
	return conn, ok
}

func (m *ActiveDBShardConnectionManager) CloseDBShards(dbID int) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	if shardMap, ok := m.ShardConns[dbID]; ok {
		for _, conn := range shardMap {
			conn.Close()
		}
		delete(m.ShardConns, dbID)
	}
}

func (m *ActiveDBShardConnectionManager) GetShardConnection(DbName string, DbID int) error {
	m.CloseDBShards(DbID)

	query := fmt.Sprintf("SELECT shard_id, shard_hash, shard_host, shard_port, shard_user, shard_pass FROM %s", DbName)
	rows, err := config.AppDBComm.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query shard mapping for DB %s: %w", DbName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var shard config.Shard
		if err := rows.Scan(&shard.ShardID, &shard.ShardHash, &shard.ShardHost, &shard.ShardPort, &shard.ShardUser, &shard.ShardPass); err != nil {
			return fmt.Errorf("failed to scan shard row: %w", err)
		}

		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", shard.ShardUser, shard.ShardPass, shard.ShardHost, shard.ShardPort, shard.ShardName)
		conn, err := sql.Open("mysql", connStr)
		if err != nil {
			return fmt.Errorf("failed to open connection for shard %d: %w", shard.ShardID, err)
		}

		if err := conn.Ping(); err != nil {
			return fmt.Errorf("failed to ping shard %d: %w", shard.ShardID, err)
		}

		m.SetShardConn(DbID, shard.ShardID, conn)
		fmt.Printf("Connected to shard %d for DB %s\n", shard.ShardID, DbName)
	}

	m.ActiveDB = DbID
	return nil
}
