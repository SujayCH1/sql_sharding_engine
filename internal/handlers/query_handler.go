package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/internal/parser"
	"sql_sharding_engine/pkg/logger"

	"github.com/xwb1989/sqlparser"
)

func HandleQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var reqBody config.Query
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := sqlparser.Parse(reqBody.QueryString)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse query: %v", err), http.StatusBadRequest)
		return
	}

	// Determine the statement type using type switch
	var handler parser.StatementHandler

	switch stmt.(type) {
	case *sqlparser.Insert:
		handler = parser.InsertHandler{}
	case *sqlparser.Select:
		handler = parser.SelectHandler{}
	case *sqlparser.Update:
		handler = parser.UpdateHandler{}
	case *sqlparser.Delete:
		handler = parser.DeleteHandler{}
	case *sqlparser.DDL:
		handler = parser.DDLHandler{}
	default:
		http.Error(w, "Unsupported statement type", http.StatusBadRequest)
		return
	}
	keys, err := handler.ExtractKeys(stmt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract primary key: %v", err), http.StatusBadRequest)
		return
	}

	logger.Logger.Info("Query received", "query", reqBody.QueryString, "primaryKeys", keys)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Primary keys: %v", keys)
}
