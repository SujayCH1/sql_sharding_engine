package cache

import (
	"context"
	"sql_sharding_engine/config"
)

func SetDBCache(ctx context.Context, name string) error {
	err := config.Redis.Set(ctx, "currDB", name, 0).Err()
	if err != nil {
		return err
	}

	config.Logger.Info("Current Dataabse set to: ", name)

	return nil
}
