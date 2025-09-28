package parser

import (
	"fmt"
	"sql_sharding_engine/internal/config"
	"strings"

	"github.com/xwb1989/sqlparser"
)

// handler to extract keys for INSERT statements.
func (h InsertHandler) ExtractKeys(stmt sqlparser.Statement) ([]string, error) {
	insertStmt, ok := stmt.(*sqlparser.Insert)
	if !ok {
		return nil, fmt.Errorf("not an insert statement")
	}

	pkIndex := -1
	for i, col := range insertStmt.Columns {
		if strings.EqualFold(col.String(), config.KeyColumn) {
			pkIndex = i
			break
		}
	}
	if pkIndex == -1 {
		return nil, fmt.Errorf("primary key column not found in INSERT")
	}

	var keys []string
	if rows, ok := insertStmt.Rows.(sqlparser.Values); ok {
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

// handler to extract keys for SELECT statements.
func (h SelectHandler) ExtractKeys(stmt sqlparser.Statement) ([]string, error) {
	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		return nil, fmt.Errorf("not a select statement")
	}
	if pkValue := extractPKFromWhere(selectStmt.Where); pkValue != "" {
		return []string{pkValue}, nil
	}
	return nil, fmt.Errorf("primary key not found in SELECT WHERE clause")
}

// handler to extract keys for UPDATE statements.
func (h UpdateHandler) ExtractKeys(stmt sqlparser.Statement) ([]string, error) {
	updateStmt, ok := stmt.(*sqlparser.Update)
	if !ok {
		return nil, fmt.Errorf("not an update statement")
	}
	if pkValue := extractPKFromWhere(updateStmt.Where); pkValue != "" {
		return []string{pkValue}, nil
	}
	return nil, fmt.Errorf("primary key not found in UPDATE WHERE clause")
}

// handler to extract keys for DELETE statements.
func (h DeleteHandler) ExtractKeys(stmt sqlparser.Statement) ([]string, error) {
	deleteStmt, ok := stmt.(*sqlparser.Delete)
	if !ok {
		return nil, fmt.Errorf("not a delete statement")
	}
	if pkValue := extractPKFromWhere(deleteStmt.Where); pkValue != "" {
		return []string{pkValue}, nil
	}
	return nil, fmt.Errorf("primary key not found in DELETE WHERE clause")
}
