package database

import (
	"database/sql"
	"strings"
)

// QueryResult wraps both rows and exec results
type QueryResult struct {
	Rows   []map[string]interface{}
	Result sql.Result
}

// ExecuteQuery executes any SQL query (SELECT, INSERT, UPDATE, DELETE, DDL)
// Returns QueryResult containing either Rows (for SELECT) or Result (for others)
func ExecuteQuery(conn *sql.DB, query string) (*QueryResult, error) {
	queryTrim := strings.TrimSpace(strings.ToUpper(query))

	if strings.HasPrefix(queryTrim, "SELECT") {
		rows, err := conn.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		cols, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		var result []map[string]interface{}

		for rows.Next() {
			values := make([]interface{}, len(cols))
			valuePtrs := make([]interface{}, len(cols))
			for i := range values {
				valuePtrs[i] = &values[i]
			}

			if err := rows.Scan(valuePtrs...); err != nil {
				return nil, err
			}

			rowMap := make(map[string]interface{})
			for i, col := range cols {
				rowMap[col] = values[i]
			}

			result = append(result, rowMap)
		}

		return &QueryResult{Rows: result}, nil
	}

	// For INSERT, UPDATE, DELETE, DDL
	res, err := conn.Exec(query)
	if err != nil {
		return nil, err
	}

	return &QueryResult{Result: res}, nil
}
