package parser

import (
	"database/sql"
	"fmt"
	"hash/crc32"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/pkg/logger"
	"strings"

	"github.com/xwb1989/sqlparser"
)

func ProcessQuery(stmt sqlparser.Statement, queryString string) error {
	if isDDL(stmt) {
		return ExecuteDDL(stmt, queryString)
	}
	return handleDML(stmt, queryString)
}

func isDDL(stmt sqlparser.Statement) bool {
	switch stmt.(type) {
	case *sqlparser.DDL:
		return true
	default:
		return false
	}
}

func handleDML(stmt sqlparser.Statement, queryString string) error {
	keys, err := extractPrimaryKeys(stmt)
	if err != nil {
		return fmt.Errorf("failed to extract primary keys: %w", err)
	}

	logger.Logger.Info("DML query received", "query", queryString, "primaryKeys", keys)

	// Execute DML with the extracted keys
	return executeDML(stmt, keys, queryString)
}

func extractPrimaryKeys(stmt sqlparser.Statement) ([]string, error) {
	switch s := stmt.(type) {
	case *sqlparser.Insert:
		return extractKeysFromInsert(s)
	case *sqlparser.Select:
		return extractKeysFromWhere(s.Where)
	case *sqlparser.Update:
		return extractKeysFromWhere(s.Where)
	case *sqlparser.Delete:
		return extractKeysFromWhere(s.Where)
	default:
		return nil, fmt.Errorf("unsupported DML statement type")
	}
}

func extractKeysFromInsert(stmt *sqlparser.Insert) ([]string, error) {
	pkIndex := -1
	for i, col := range stmt.Columns {
		if strings.EqualFold(col.String(), config.KeyColumn) {
			pkIndex = i
			break
		}
	}

	if pkIndex == -1 {
		return nil, fmt.Errorf("primary key column not found in INSERT")
	}

	var keys []string
	if rows, ok := stmt.Rows.(sqlparser.Values); ok {
		for _, row := range rows {
			if pkIndex < len(row) {
				keys = append(keys, sqlparser.String(row[pkIndex]))
			}
		}
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("no primary keys found in INSERT rows")
	}

	return keys, nil
}

func extractKeysFromWhere(where *sqlparser.Where) ([]string, error) {
	if where == nil {
		return nil, fmt.Errorf("WHERE clause is required for this operation")
	}

	comparison, ok := where.Expr.(*sqlparser.ComparisonExpr)
	if !ok {
		return nil, fmt.Errorf("unsupported WHERE expression")
	}

	if strings.EqualFold(sqlparser.String(comparison.Left), config.KeyColumn) &&
		comparison.Operator == "=" {
		pkValue := sqlparser.String(comparison.Right)
		return []string{pkValue}, nil
	}

	return nil, fmt.Errorf("primary key not found in WHERE clause")
}

// getShardIDForKey computes hash of primary key and maps it to one of the available shard IDs
func getShardIDForKey(key string, shards map[int]*sql.DB) (int, error) {
	if len(shards) == 0 {
		return 0, fmt.Errorf("no shards provided")
	}

	// Simple consistent hashing using crc32
	hash := crc32.ChecksumIEEE([]byte(key))
	shardIDs := make([]int, 0, len(shards))
	for id := range shards {
		shardIDs = append(shardIDs, id)
	}

	// Map hash to a shard
	index := int(hash) % len(shardIDs)
	return shardIDs[index], nil
}
