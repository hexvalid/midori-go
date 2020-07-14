package database

import (
	"database/sql"
	"github.com/fatih/color"
	"github.com/hexvalid/midori-go/logger"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var log = logger.NewLog("Database", color.FgCyan)

func Open(file string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		_, err = db.Exec(initQuery)
		if err != nil {
			return nil, err
		}
		log.Info("Database created and opened succesfully: %s", color.YellowString(file))
	} else {
		log.Info("Database opened succesfully: %s", color.YellowString(file))
	}
	return db, nil
}

func Close(db *sql.DB) error {
	log.Info("Closing database...")
	return db.Close()
}
