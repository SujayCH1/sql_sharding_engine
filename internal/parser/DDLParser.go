package parser

import "github.com/xwb1989/sqlparser"

// Hnadler for DLL statements
func (h DDLHandler) ExtractKeys(stmt sqlparser.Statement) ([]string, error) {
	// For DDL we don't need primary keys; just return nil
	return nil, nil
}
