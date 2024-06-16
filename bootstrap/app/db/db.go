package db

import (
	"log"
	"os"

	"github.com/anthdm/superkit/db"
	"github.com/anthdm/superkit/kit"

	_ "github.com/mattn/go-sqlite3"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
)

// I could not came up with a better naming for this.
// Ideally, app should export a global variable called "DB"
// but this will cause imports cycles for plugins.
var Query *bun.DB

func init() {
	config := db.Config{
		Driver:   os.Getenv("DB_DRIVER"),
		Name:     os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		Host:     os.Getenv("DB_HOST"),
	}
	db, err := db.New(config)
	if err != nil {
		log.Fatal(err)
	}
	Query = bun.NewDB(db, sqlitedialect.New())
	if kit.IsDevelopment() {
		Query.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
}
