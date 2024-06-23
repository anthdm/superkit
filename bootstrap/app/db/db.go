package db

import (
	"log"
	"os"

	"github.com/anthdm/superkit/db"

	_ "github.com/mattn/go-sqlite3"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// By default this is a pre-configured Gorm DB instance.
// Change this type based on the database package of your likings.
var dbInstance *gorm.DB

// Get returns the instantiated DB instance.
func Get() *gorm.DB {
	return dbInstance
}

func init() {
	// Create a default *sql.DB exposed by the superkit/db package
	// based on the given configuration.
	config := db.Config{
		Driver:   os.Getenv("DB_DRIVER"),
		Name:     os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
	}
	dbinst, err := db.NewSQL(config)
	if err != nil {
		log.Fatal(err)
	}
	// Based on the driver create the corresponding DB instance.
	// By default, the SuperKit boilerplate comes with a pre-configured
	// ORM called Gorm. https://gorm.io.
	//
	// You can change this to any other DB interaction tool
	// of your liking. EG:
	// - uptrace bun -> https://bun.uptrace.dev
	// - SQLC -> https://github.com/sqlc-dev/sqlc
	// - gojet -> https://github.com/go-jet/jet
	switch config.Driver {
	case db.DriverSqlite3:
		dbInstance, err = gorm.Open(sqlite.New(sqlite.Config{
			Conn: dbinst,
		}))
	case db.DriverMysql:
		// ...
	default:
		log.Fatal("invalid driver:", config.Driver)
	}
	if err != nil {
		log.Fatal(err)
	}
}
