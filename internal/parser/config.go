package parser

import "github.com/xwb1989/sqlparser"

// StatementHandler defines interface for all statement types
type StatementHandler interface {
	ExtractKeys(stmt sqlparser.Statement) ([]string, error)
}

// Mapping statement types to handlers
var handlers = map[string]StatementHandler{
	"insert": InsertHandler{},
	"select": SelectHandler{},
	"update": UpdateHandler{},
	"delete": DeleteHandler{},
	"ddl":    DDLHandler{},
}

// InsertHandler handles INSERT statements
type InsertHandler struct{}

// SelectHandler handles SELECT statements
type SelectHandler struct{}

// UpdateHandler handles UPDATE statements
type UpdateHandler struct{}

// DeleteHandler handles DELETE statements
type DeleteHandler struct{}

// DDLHandler handles CREATE, ALTER, DROP, etc.
type DDLHandler struct{}
