package main

import (
	"sql_sharding_engine/internal/loader"
	"sql_sharding_engine/pkg/logger"
)

func main() {
	err := loader.LoadServices()
	if err != nil {
		logger.Logger.Error("failed to load application services:", "error", err)
	}

	select {}

}
