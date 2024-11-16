// This file holds the SQLite database connection handler
package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Vikuuu/gobase/internal/generator"
)

const DBFILENAME = "db.sqlite3"

func SqliteConn(dbFileName string) (*sql.DB, error) {
	db, err := sql.Open(SQLITE, dbFileName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB, fileName string) error {
	query := generator.SqLiteCreateTable(fileName)

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
