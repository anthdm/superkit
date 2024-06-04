package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DriverSqlite3 = "sqlite3"
)

func New() (*sql.DB, error) {
	driver := os.Getenv("DB_DRIVER")

	switch driver {
	case DriverSqlite3:
		name := os.Getenv("DB_NAME")
		if len(name) == 0 {
			name = "gothkit"
		}
		return sql.Open(driver, name)
	default:
		return nil, fmt.Errorf("invalid database driver (%s): currently only sqlite3 is supported", driver)
	}
}
