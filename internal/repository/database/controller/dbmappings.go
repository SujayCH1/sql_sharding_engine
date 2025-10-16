package database

import (
	"fmt"
	"sql_sharding_engine/internal/config"
)

func FetchDBMappings() ([]RowsData, error) {
	// Check DB connection
	if config.AppDBComm == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	query := fmt.Sprintf("SELECT * FROM dbmappings")

	rows, err := config.AppDBComm.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query DB mapping for view: %w", err)
	}
	defer rows.Close()

	var data []RowsData

	for rows.Next() {
		var temp RowsData
		if err := rows.Scan(&temp.ID, &temp.Name); err != nil {
			return nil, fmt.Errorf("failed to parse rows for view: %w", err)
		}
		data = append(data, temp)
	}

	return data, nil
}
