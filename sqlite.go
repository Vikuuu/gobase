package gobase

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const DBFILENAME = "db.sqlite3"

func SqliteConn(dbFileName string) (*sql.DB, error) {
	db, err := sql.Open(SQLITE, dbFileName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
