package connections

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sql_sharding_engine/internal/config"
	"sql_sharding_engine/pkg/logger"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// func to estavlish connection with app DB and store in global variable
func LoadMainDBConn() error {

	err := AddMainDBCred()
	if err != nil {
		return fmt.Errorf("failed to add app DB credentials: %w", err)
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.AppDBCommInfo.DBUser,
		config.AppDBCommInfo.DBPass,
		config.AppDBCommInfo.DBHost,
		config.AppDBCommInfo.DBPort,
		config.AppDBCommInfo.DBName,
	)

	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		return fmt.Errorf("failed to load app DB connection: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	config.AppDBComm = conn

	logger.Logger.Info("Established connection with application database.")

	return nil
}

// func to establish connection with shards
func LoadDBConn(c config.DBConnInfo) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)

	conn, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to load shard connection: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping shard: %w", err)
	}

	logger.Logger.Info("Established connection with shard: ", c.DBName, ".")

	return conn, err
}

// func to establish and store connection detials and credentails of  application database
func AddMainDBCred() error {

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return errors.New("unable to convert port string")
	}

	cfg := config.DBConnInfo{
		DBName: os.Getenv("DB_NAME"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: port,
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
	}

	config.AppDBCommInfo = &cfg

	return nil

}
