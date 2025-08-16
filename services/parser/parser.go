package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sql_sharding_engine/config"
	"strings"

	"github.com/xwb1989/sqlparser"
)

// Handler function for /query route
// Parese query and timestamp for req body
// Extracts primary key from query
// Calls hasher to find expected shard credentials
func HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody config.Query

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	primaryKey, err := extractKeys(reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract primary key: %v", err), http.StatusBadRequest)
		return
	}

	config.Logger.Info("Primary key extracted")

	w.WriteHeader(http.StatusOK)
	config.Logger.Info("Query received", "query", reqBody.QueryString, "primaryKey", primaryKey)

}

// helper to parse and find pk from query string
func extractKeys(q config.Query) ([]string, error) {
	stmt, err := sqlparser.Parse(q.QueryString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %w", err)
	}

	switch stmt := stmt.(type) {
	case *sqlparser.Insert:
		pkIndex := -1
		for i, col := range stmt.Columns {
			if strings.EqualFold(col.String(), "primary_key") {
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

	case *sqlparser.Update:
		if pkValue := extractPKFromExpr(stmt.Where); pkValue != "" {
			return []string{pkValue}, nil
		}
		return nil, fmt.Errorf("primary key not found in UPDATE WHERE clause")

	case *sqlparser.Delete:
		if pkValue := extractPKFromExpr(stmt.Where); pkValue != "" {
			return []string{pkValue}, nil
		}
		return nil, fmt.Errorf("primary key not found in DELETE WHERE clause")

	case *sqlparser.Select:
		if pkValue := extractPKFromExpr(stmt.Where); pkValue != "" {
			return []string{pkValue}, nil
		}
		return nil, fmt.Errorf("primary key not found in SELECT WHERE clause")
	}

	return nil, fmt.Errorf("unsupported statement type")
}

// helper to extact pk from SQL expression
func extractPKFromExpr(where *sqlparser.Where) string {
	if where == nil {
		return ""
	}
	comparison, ok := where.Expr.(*sqlparser.ComparisonExpr)
	if !ok {
		return ""
	}
	if strings.EqualFold(sqlparser.String(comparison.Left), "primary_key") &&
		comparison.Operator == "=" {
		return sqlparser.String(comparison.Right)
	}
	return ""
}
