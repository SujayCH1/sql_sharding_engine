package cache

import (
	"context"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/pkg/logger"
)

func SetDBCache(ctx context.Context, name string) error {
	err := config.Redis.Set(ctx, "currDB", name, 0).Err()
	if err != nil {
		return err
	}

	logger.Logger.Info("Current Dataabse set to: ", name)

	return nil
}
