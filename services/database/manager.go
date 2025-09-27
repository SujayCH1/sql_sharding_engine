package database

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sql_sharding_engine/config"
)

// func to add/remove database
func HandleDatabase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newReq dbReq

	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	switch newReq.ReqType {
	case "add":
		err := AddDatabase(newReq.DBInfo)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			config.Logger.Error("failed to add a new database", err)
		}

	case "delete":
		err := DeleteDatabase(newReq.DBInfo)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			config.Logger.Error("failed to delete exsisting database:", err)
		}

	default:
		http.Error(w, "invalid type", http.StatusBadRequest)
	}

}

// func to add db mapping
func AddDatabase(d config.Database) error {
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

	config.Logger.Info("Database added:", name)
	return nil
}

// fun to remove db mappings
func DeleteDatabase(d config.Database) error {
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

	config.Logger.Info("Database deleted:", name)
	return nil
}

// func to select db
func (m *CurrDBManager) HandleSelectDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var db CurrDB
	err := json.NewDecoder(r.Body).Decode(&db)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	m.SetCurrentDB(&db)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Current DB set to: %s", db.Name)

	config.Logger.Info("Current DB updated", "db", db.Name)
}
