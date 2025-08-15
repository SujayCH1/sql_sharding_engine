package parser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sql_sharding_engine/services"

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

	services.Logger.Info("Query reqest accpeted")

	var reqBody services.Query

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	services.Logger.Info("Query request parsed")

	primaryKey, err := extractKey(reqBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract primary key: %v", err), http.StatusBadRequest)
		return
	}

	services.Logger.Info("Primary key extracted")

	w.WriteHeader(http.StatusOK)
	services.Logger.Info("Query received", "query", reqBody.QueryString, "primaryKey", primaryKey)

}

// helper to parse and find pk from query string
func extractKey(q services.Query) (string, error) {
	stmt, err := sqlparser.Parse(q.QueryString)
	if err != nil {
		return "", fmt.Errorf("failed to parse statement: %s", err)
	}

	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		if stmt.Where != nil {
			return extractPKFromExpr(stmt.Where.Expr, services.KeyColumn)
		}

	case *sqlparser.Insert:
		cols := stmt.Columns
		for i, col := range cols {
			if col.String() == services.KeyColumn {
				if rows, ok := stmt.Rows.(sqlparser.Values); ok && len(rows) > 0 && len(rows[0]) > i {
					return sqlparser.String(rows[0][i]), nil
				}
			}
		}

	case *sqlparser.Update:
		if stmt.Where != nil {
			return extractPKFromExpr(stmt.Where.Expr, services.KeyColumn)
		}

	case *sqlparser.Delete:
		if stmt.Where != nil {
			return extractPKFromExpr(stmt.Where.Expr, services.KeyColumn)
		}
	}

	return "", fmt.Errorf("unsupported statement type or primary key not found")
}

// helper to extact pk from SQL expression
func extractPKFromExpr(expr sqlparser.Expr, pkColumn string) (string, error) {
	if comp, ok := expr.(*sqlparser.ComparisonExpr); ok {
		if colName, ok := comp.Left.(*sqlparser.ColName); ok {
			if colName.Name.String() == pkColumn {
				return sqlparser.String(comp.Right), nil
			}
		}
	}

	if andExpr, ok := expr.(*sqlparser.AndExpr); ok {
		if val, err := extractPKFromExpr(andExpr.Left, pkColumn); err == nil {
			return val, nil
		}
		return extractPKFromExpr(andExpr.Right, pkColumn)
	}

	return "", fmt.Errorf("primary key column not found in expression")
}
