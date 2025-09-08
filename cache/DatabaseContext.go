package cache

import (
	"context"
	"encoding/json"
	"net/http"
	"sql_sharding_engine/config"
)

type currDB struct {
	Name string `json:"name`
	ID   int    `json:"id`
}

func HandleSelectDB(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var db currDB

	err := json.NewDecoder(r.Body).Decode(&db)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	err = setDBCache(context.Background(), db.Name)
	if err != nil {
		http.Error(w, "Failed to add DB Cache", http.StatusInternalServerError)
		config.Logger.Error("Failed to add DB Cache:", err)
		return
	}

}

func setDBCache(ctx context.Context, name string) error {
	err := config.Redis.Set(ctx, "currDB", name, 0).Err()
	if err != nil {
		return err
	}

	config.Logger.Info("Current Dataabse set to: ", name)

	return nil
}
