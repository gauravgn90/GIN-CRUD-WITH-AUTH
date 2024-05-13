package connection

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB(dataSourceName string, maxOpenConns, maxIdleConns int) error {
	// Open a new database connection with a connection pool
	conn, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	// Set maximum open and idle connections
	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetMaxIdleConns(maxIdleConns)

	// Test the connection
	if err := conn.Ping(); err != nil {
		return err
	}

	db = conn
	return nil
}

func GetDB() *sql.DB {
	// Return the existing database connection
	return db
}
