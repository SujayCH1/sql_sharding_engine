package main

import (
	"sql_sharding_engine/config"
	"sql_sharding_engine/loader"
)

func main() {
	err := loader.LoadServices()
	if err != nil {
		config.Logger.Error("failed to load application services:", "error", err)
	}

	select {}

}
