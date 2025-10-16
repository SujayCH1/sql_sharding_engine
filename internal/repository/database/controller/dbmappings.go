package database

import (
	"fmt"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/internal/repository/database"
)

func FetchDBMappings() ([]RowsData, error) {

	query := fmt.Sprintf("SELECT * FROM %s", database.CurrDBMgr.GetCurrentDB().Name)

	rows, err := config.AppDBComm.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query DB mapping for view  %s: %w", config.AppDBCommInfo.DBName, err)
	}

	defer rows.Close()

	var Data []RowsData

	for rows.Next() {
		var temp RowsData

		if err := rows.Scan(&temp.ID, &temp.Name); err != nil {
			return nil, fmt.Errorf("failed to parrse rows for view : %s: %w", config.AppDBCommInfo.DBName, err)
		}

		Data = append(Data, temp)
	}

	return Data, nil
}
