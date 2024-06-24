package db

import (
	"database/sql"
	"fmt"
)

const (
	DriverSqlite3 = "sqlite3"
	DriverMysql   = "mysql"
)

type Config struct {
	Driver   string
	Name     string
	Host     string
	User     string
	Password string
}

func NewSQL(cfg Config) (*sql.DB, error) {
	switch cfg.Driver {
	case DriverSqlite3:
		name := cfg.Name
		if len(name) == 0 {
			name = "app_db"
		}
		return sql.Open(cfg.Driver, name)
	default:
		return nil, fmt.Errorf("invalid database driver (%s): currently only sqlite3 is supported", cfg.Driver)
	}
}
