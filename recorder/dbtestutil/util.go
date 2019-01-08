package dbtestutil

import (
	"database/sql"
	"fmt"
)

var conn *sql.DB

// DB constants for testing
const (
	DBUsername = "root"
	DBPassword = "my-secret"
	DBHost     = "localhost"
	DBPort     = "3308"
	DBName     = "lino_db"
)

// NewDBConn returns a new sql db connection
func NewDBConn() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		DBUsername,
		DBPassword,
		DBHost,
		DBPort,
		DBName,
	)
	conn, err := sql.Open("mysql", dsn)
	return conn, err
}
