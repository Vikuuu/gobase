package database

import (
	"database/sql"
	"errors"
	"os"
	"testing"
)

func TestSqliteConn(t *testing.T) {
	db, err := SqliteConn(DBFILENAME)
	if err != nil {
		t.Fatalf("Error creating conn. error=%s", err)
	}

	db.Close()
}

func TestCreateTable(t *testing.T) {
	fileName := "../test/create_table.go"

	// Creating conn with DB
	db, err := SqliteConn(DBFILENAME)
	if err != nil {
		t.Fatalf("Error creating conn. error=%s", err)
	}
	defer db.Close()

	// Creating Table in the conn
	err = createTable(db, fileName)
	if err != nil {
		t.Fatalf("Error creating table. error=%s", err)
	}

	// Check if the table exists
	row := db.QueryRow(
		`SELECT name FROM sqlite_schema WHERE type='table' AND name NOT LIKE 'sqlite_%'`,
	)
	if err != nil {
		t.Fatalf("Error executing stmt. error=%s", err)
	}

	var output struct {
		got string
	}
	if err := row.Scan(&output.got); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("No row found. error=%s", err)
		}
		t.Fatalf("Error scanning. error=%s", err)
	}

	if "users" != output.got {
		t.Fatalf("Table not equals. expected=users. got=%s", output.got)
	}

	// Remove the file after testing
	os.Remove(DBFILENAME)
}
