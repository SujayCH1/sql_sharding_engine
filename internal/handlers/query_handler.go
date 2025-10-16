package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/internal/parser"

	"github.com/xwb1989/sqlparser"
)

// HandleQuery is the main HTTP handler for query processing
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

	// Process the query based on its type
	if err := parser.ProcessQuery(stmt, reqBody.QueryString); err != nil {
		http.Error(w, fmt.Sprintf("Failed to process query: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Query processed successfully")
}
