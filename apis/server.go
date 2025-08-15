package apis

import (
	"net/http"
	"sql_sharding_engine/services"
	"sql_sharding_engine/services/parser"
	"time"
)

func StartServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/query", parser.HandleQuery)

	server := &http.Server{
		Addr:         ":8085",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	services.Logger.Info("Server listening at port 8085.....")
	server.ListenAndServe()
}
