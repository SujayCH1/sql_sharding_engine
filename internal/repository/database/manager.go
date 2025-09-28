package database

import (
	"fmt"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/pkg/logger"
)

// func to add db mapping
func AddDatabase(d Database) error {
	name := d.Name
	id := d.ID

	query1 := `
		INSERT INTO dbmappings (id, name) VALUES (?, ?)
	`

	query2 := fmt.Sprintf(`
		CREATE TABLE %s (
			database_name VARCHAR(50) NOT NULL,
			shard_id INT NOT NULL PRIMARY KEY,
			shard_hash VARCHAR(50) NOT NULL,
			shard_host VARCHAR(50) NOT NULL,
			shard_port INT NOT NULL,
			shard_user VARCHAR(50) NOT NULL,
			shard_pass VARCHAR(50) NOT NULL
		);`, name)

	_, err := config.AppDBComm.Exec(query1, id, name)
	if err != nil {
		return fmt.Errorf("failed to insert into dbmappings: %w", err)
	}

	_, err = config.AppDBComm.Exec(query2)
	if err != nil {
		return fmt.Errorf("failed to create database %s: %w", name, err)
	}

	logger.Logger.Info("Database added:", name)
	return nil
}

// fun to remove db mappings
func DeleteDatabase(d Database) error {
	name := d.Name
	id := d.ID

	query1 := `DELETE FROM dbmappings WHERE id = ? AND name = ?`

	query2 := fmt.Sprintf(`DROP TABLE %s`, name)

	_, err := config.AppDBComm.Exec(query1, id, name)
	if err != nil {
		return fmt.Errorf("failed to delete mapping for %s: %w", name, err)
	}

	_, err = config.AppDBComm.Exec(query2)
	if err != nil {
		return fmt.Errorf("failed to drop table %s: %w", name, err)
	}

	logger.Logger.Info("Database deleted:", name)
	return nil
}
